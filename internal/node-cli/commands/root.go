// Copyright © 2017 The virtual-kubelet authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"os"
	"path"
	"sync"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/manager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/opts"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/provider"

	"github.com/pkg/errors"
	"github.com/virtual-kubelet/virtual-kubelet/errdefs"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	"github.com/virtual-kubelet/virtual-kubelet/node"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes/typed/coordination/v1beta1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"
)

var sharedFactory kubeinformers.SharedInformerFactory
var once sync.Once

// MetricsPortKey determines the key of the value in the context, which can be used to retrieve the metrics port
const MetricsPortKey = "metrics-port"

func RunRootCommand(originalCtx context.Context, ctx context.Context, s *provider.Store, c *opts.Opts) (int, int, error) {
	pInit := s.Get(c.Provider)
	if pInit == nil {
		return 0, 0, errors.Errorf("provider %q not found", c.Provider)
	}

	client, err := newClient(c.KubeConfigPath, c.KubeAPIQPS, c.KubeAPIBurst)
	if err != nil {
		return 0, 0, err
	}

	return runRootCommandWithProviderAndClient(originalCtx, ctx, pInit, client, c)
}

func runRootCommandWithProviderAndClient(originalCtx context.Context, ctx context.Context, pInit provider.InitFunc, client kubernetes.Interface, c *opts.Opts) (int, int, error) {
	if ok := provider.ValidOperatingSystems[c.OperatingSystem]; !ok {
		return 0, 0, errdefs.InvalidInputf("operating system %q is not supported", c.OperatingSystem)
	}

	if c.PodSyncWorkers == 0 {
		return 0, 0, errdefs.InvalidInput("pod sync workers must be greater than 0")
	}

	var taint *corev1.Taint
	if !c.DisableTaint {
		var err error
		taint, err = getTaint(c)
		if err != nil {
			return 0, 0, err
		}
	}

	// Create a shared informer factory for Kubernetes pods in the current namespace (if specified) and scheduled to the current node.
	podInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(
		client,
		c.InformerResyncPeriod,
		kubeinformers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.OneTermEqualSelector("spec.nodeName", c.NodeName).String()
		}))
	podInformer := podInformerFactory.Core().V1().Pods()

	// Create another shared informer factory for Kubernetes secrets and configmaps (not subject to any selectors).
	once.Do(func() {
		sharedFactory = kubeinformers.NewSharedInformerFactoryWithOptions(client, c.InformerResyncPeriod)
	})

	// Create a secret informer and a config map informer so we can pass their listers to the resource manager.
	secretInformer := sharedFactory.Core().V1().Secrets()
	configMapInformer := sharedFactory.Core().V1().ConfigMaps()
	serviceInformer := sharedFactory.Core().V1().Services()

	rm, err := manager.NewResourceManager(podInformer.Lister(), secretInformer.Lister(), configMapInformer.Lister(), serviceInformer.Lister())
	if err != nil {
		return 0, 0, errors.Wrap(err, "could not create resource manager")
	}

	// Start the informers now, so the provider will get a functional resource
	// manager.
	sharedFactory.Start(originalCtx.Done())
	podInformerFactory.Start(ctx.Done())

	apiConfig, err := getAPIConfig(c)
	if err != nil {
		return 0, 0, err
	}

	initConfig := provider.InitConfig{
		ConfigPath:        c.ProviderConfigPath,
		NodeName:          c.NodeName,
		OperatingSystem:   c.OperatingSystem,
		ResourceManager:   rm,
		DaemonPort:        int32(c.ListenPort),
		InternalIP:        os.Getenv("VKUBELET_POD_IP"),
		KubeClusterDomain: c.KubeClusterDomain,
	}

	p, err := pInit(&initConfig)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "error initialising provider %s", c.Provider)
	}

	cancelHTTP, metricsPort, k8sPort, err := setupHTTPServer(ctx, p, apiConfig)
	if err != nil {
		return 0, 0, err
	}
	c.ListenPort = int32(k8sPort)
	initConfig.DaemonPort = int32(k8sPort)
	ctx = context.WithValue(ctx, MetricsPortKey, metricsPort)

	go func() {
		defer cancelHTTP()

		<-ctx.Done()
	}()

	ctx = log.WithLogger(ctx, log.G(ctx).WithFields(log.Fields{
		"provider":         c.Provider,
		"operatingSystem":  c.OperatingSystem,
		"node":             c.NodeName,
		"watchedNamespace": c.KubeNamespace,
	}))

	var leaseClient v1beta1.LeaseInterface
	if c.EnableNodeLease {
		leaseClient = client.CoordinationV1beta1().Leases(corev1.NamespaceNodeLease)
	}

	nodeProvider, ok := p.(node.NodeProvider)
	if !ok {
		nodeProvider = node.NaiveNodeProvider{}
	}
	pNode := NodeFromProvider(ctx, c.NodeName, taint, p, c.Version)
	nodeRunner, err := node.NewNodeController(
		nodeProvider,
		pNode,
		client.CoreV1().Nodes(),
		node.WithNodeEnableLeaseV1Beta1(leaseClient, nil),
		node.WithNodeStatusUpdateErrorHandler(func(ctx context.Context, err error) error {
			if !k8serrors.IsNotFound(err) {
				return err
			}

			log.G(ctx).Debug("node not found")
			newNode := pNode.DeepCopy()
			newNode.ResourceVersion = ""
			_, err = client.CoreV1().Nodes().Create(newNode)
			if err != nil {
				return err
			}
			log.G(ctx).Debug("created new node")
			return nil
		}),
	)
	if err != nil {
		return 0, 0, errors.Wrap(err, "unable to create node controller")
	}

	eb := record.NewBroadcaster()
	eb.StartLogging(log.G(ctx).Infof)
	eb.StartRecordingToSink(&corev1client.EventSinkImpl{Interface: client.CoreV1().Events(c.KubeNamespace)})

	pc, err := node.NewPodController(node.PodControllerConfig{
		PodClient:         client.CoreV1(),
		PodInformer:       podInformer,
		EventRecorder:     eb.NewRecorder(scheme.Scheme, corev1.EventSource{Component: path.Join(pNode.Name, "pod-controller")}),
		Provider:          p,
		SecretInformer:    secretInformer,
		ConfigMapInformer: configMapInformer,
		ServiceInformer:   serviceInformer,
	})
	if err != nil {
		return 0, 0, errors.Wrap(err, "error setting up pod controller")
	}

	go func() {
		if err := pc.Run(ctx, c.PodSyncWorkers); err != nil && errors.Cause(err) != context.Canceled {
			log.G(ctx).Fatal(err)
		}
	}()

	if c.StartupTimeout > 0 {
		// If there is a startup timeout, it does two things:
		// 1. It causes the VirtualKubelet to shutdown if we haven't gotten into an operational state in a time period
		// 2. It prevents node advertisement from happening until we're in an operational state
		err = waitFor(ctx, c.StartupTimeout, pc.Ready())
		if err != nil {
			return 0, 0, err
		}
	}

	go func() {
		if err := nodeRunner.Run(ctx); err != nil {
			log.G(ctx).Fatal(err)
		}
	}()

	log.G(ctx).Info("Initialised")

	return metricsPort, k8sPort, nil
}

func waitFor(ctx context.Context, time time.Duration, ready <-chan struct{}) error {
	ctx, cancel := context.WithTimeout(ctx, time)
	defer cancel()

	// Wait for the VirtualKubelet / PC close the the ready channel, or time out and return
	log.G(ctx).Info("Waiting for pod controller / VirtualKubelet to be ready")

	select {
	case <-ready:
		return nil
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "Error while starting up VirtualKubelet")
	}
}

func newClient(configPath string, qps, burst int32) (*kubernetes.Clientset, error) {
	var config *rest.Config

	// Check if the kubeConfig file exists.
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		// Get the kubeconfig from the filepath.
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, errors.Wrap(err, "error building client config")
		}
	} else {
		// Set to in-cluster config.
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, errors.Wrap(err, "error building in cluster config")
		}
	}

	if qps != 0 {
		config.QPS = float32(qps)
	}

	if burst != 0 {
		config.Burst = int(burst)
	}

	if masterURI := os.Getenv("MASTER_URI"); masterURI != "" {
		config.Host = masterURI
	}

	return kubernetes.NewForConfig(config)
}
