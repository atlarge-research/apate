package provider

import (
	"testing"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/provider"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

const podNamespace = "podnamespace"
const podName = "pod"
const podLabel = "label"

func TestNewProvider(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	pm := podmanager.New()
	sts := NewStats()
	resources := scenario.NodeResources{
		UUID:             uuid.New(),
		Memory:           0,
		CPU:              0,
		Storage:          0,
		EphemeralStorage: 0,
		MaxPods:          0,
		Label:            "",
	}

	cfg := provider.InitConfig{}
	ni, err := node.NewInfo("a", "b", "c", "d", "e/f")
	assert.NoError(t, err)

	var s store.Store = ms

	ms.EXPECT().AddPodFlagListener(events.PodResources, gomock.Any())

	e, err := env.ApateletEnv()
	assert.NoError(t, err)
	p, ok := NewProvider(pm, sts, &resources, &cfg, &ni, &s, true, e).(*Provider)

	assert.True(t, ok)

	assert.EqualValues(t, p.Conditions.ready.Get().Status, metav1.ConditionTrue)
	assert.EqualValues(t, p.Conditions.outOfDisk.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.memoryPressure.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.diskPressure.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.networkUnavailable.Get().Status, metav1.ConditionFalse)
	assert.EqualValues(t, p.Conditions.pidPressure.Get().Status, metav1.ConditionFalse)
}
