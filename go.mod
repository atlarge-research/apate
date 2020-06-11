module github.com/atlarge-research/opendc-emulate-kubernetes

go 1.14

require (
	github.com/deanishe/go-env v0.4.0
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.3.0
	github.com/docker/go-units v0.3.3
	github.com/fatih/color v1.9.0
	github.com/finitum/node-cli v0.1.3-0.20200611095742-0bf9cf7cee8e
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli/v2 v2.2.0
	github.com/virtual-kubelet/virtual-kubelet v1.2.1
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.28.1
	google.golang.org/protobuf v1.21.0
	k8s.io/api v0.18.2 // Will be replaced
	k8s.io/apimachinery v0.18.2 // Will be replaced
	k8s.io/client-go v10.0.0+incompatible // Will be replaced
	k8s.io/kubernetes v1.15.2
	sigs.k8s.io/kind v0.7.0
)

// Replaced due to security bug: https://nvd.nist.gov/vuln/detail/CVE-2019-0210
replace github.com/apache/thrift v0.12.0 => github.com/apache/thrift v0.13.0

// Replaced due to security bug: https://access.redhat.com/errata/RHBA-2018:2652
replace github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.5.0+incompatible

// Replaced due to security bug: https://github.com/docker/cli/pull/2117
replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

// TODO: Figure out why we want all these replacements
//       taken from: https://github.com/virtual-kubelet/virtual-kubelet/blob/master/go.mod#L48

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190805141119-fdd30b57c827
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190805143126-cdb999c96590
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190805142138-368b2058237c
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190805143448-a07e59fb081d
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190805141520-2fe0317bcee0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190805144409-8484242760e7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190805144246-c01ee70854a1
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190805141645-3a5e5ac800ae
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190531030430-6117653b35f1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190805144531-3985229e1802
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190805142416-fd821fbbb94e
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190805144128-269742da31dd
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190805143734-7f1675b90353
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190805144012-2a1ed1f3d8a4
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190805143852-517ff267f8d1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20190805144654-3d5bf3a310c1
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190805143318-16b07057415d
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190805142637-3b65bc4bb24f
)
