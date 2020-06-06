package pod

import (
	"testing"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/docker/go-units"
	stats "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	podconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store/mock_store"
)

func TestGetCRDAndLabel(t *testing.T) {
	t.Parallel()

	ep := podconfigv1.PodConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
	}

	lbl := getCrdLabel(&ep)

	assert.Equal(t, lbl, "TestNamespace/TestName")
}

func TestEnqueueCRD(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	et1 := store.NewPodTask(
		1*time.Millisecond,
		"TestNamespace/TestName", &podconfigv1.PodConfigurationState{
			PodStatus: podconfigv1.PodStatusFailed,
		})

	et2 := store.NewPodTask(
		42*time.Millisecond,
		"TestNamespace/TestName", &podconfigv1.PodConfigurationState{
			PodStatus: podconfigv1.PodStatusPending,
		})

	ep := podconfigv1.PodConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: podconfigv1.PodConfigurationSpec{
			Tasks: []podconfigv1.PodConfigurationTask{
				{
					Timestamp: "1ms",
					State: podconfigv1.PodConfigurationState{
						PodStatus: podconfigv1.PodStatusFailed,
					},
				},
				{
					Timestamp: "42ms",
					State: podconfigv1.PodConfigurationState{
						PodStatus: podconfigv1.PodStatusPending,
					},
				},
				{
					Timestamp:     "1s",
					RelativeToPod: true,
					State: podconfigv1.PodConfigurationState{
						PodStatus: podconfigv1.PodStatusRunning,
					},
				},
			},
		},
	}

	ms.EXPECT().SetPodTasks(
		"TestNamespace/TestName",
		gomock.Any(),
	).Do(func(_ string, arr []*store.Task) {
		assert.Equal(t, 2, len(arr))
		assert.EqualValues(t, et1, arr[0])
		assert.EqualValues(t, et2, arr[1])
	})

	ms.EXPECT().SetPodTimeFlags(
		"TestNamespace/TestName",
		gomock.Any(),
	).Do(func(_ string, arr []*store.TimeFlags) {
		assert.Equal(t, 1, len(arr))
		assert.EqualValues(t, &store.TimeFlags{
			TimeSincePodStart: 1 * time.Second,
			Flags: store.Flags{
				events.PodStatus: scenario.PodStatusRunning,
			},
		}, arr[0])
	})

	err := setPodTasks(&ep, &s)
	assert.NoError(t, err)
}

func TestEnqueueCRDDirect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ms := mock_store.NewMockStore(ctrl)

	var s store.Store = ms

	ep := podconfigv1.PodConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "TestName",
			Namespace: "TestNamespace",
		},
		Spec: podconfigv1.PodConfigurationSpec{
			PodConfigurationState: podconfigv1.PodConfigurationState{
				CreatePodResponse:    podconfigv1.ResponseNormal,
				UpdatePodResponse:    podconfigv1.ResponseNormal,
				DeletePodResponse:    podconfigv1.ResponseNormal,
				GetPodResponse:       podconfigv1.ResponseNormal,
				GetPodStatusResponse: podconfigv1.ResponseNormal,
				PodResources: &podconfigv1.PodResources{
					Memory:           "10T",
					CPU:              1000,
					Storage:          "5K",
					EphemeralStorage: "100M",
				},
				PodStatus: podconfigv1.PodStatusRunning,
			},
			Tasks: []podconfigv1.PodConfigurationTask{},
		},
	}

	cores := uint64(1000)
	memory := uint64(10 * units.TiB)
	storage := uint64(5 * units.KiB)
	ephStorage := uint64(100 * units.MiB)

	ms.EXPECT().SetPodFlags("TestNamespace/TestName", gomock.Any()).Do(func(_ string, flags store.Flags) {
		assert.Equal(t, events.PodCreatePodResponse, translateResponse(podconfigv1.ResponseNormal))
		assert.Equal(t, events.PodUpdatePodResponse, translateResponse(podconfigv1.ResponseNormal))
		assert.Equal(t, events.PodDeletePodResponse, translateResponse(podconfigv1.ResponseNormal))
		assert.Equal(t, events.PodGetPodResponse, translateResponse(podconfigv1.ResponseNormal))
		assert.Equal(t, events.PodGetPodStatusResponse, translateResponse(podconfigv1.ResponseNormal))

		stat := flags[events.PodResources].(*stats.PodStats)

		assert.EqualValues(t, cores, *stat.CPU.UsageNanoCores)
		assert.WithinDuration(t, time.Now(), stat.CPU.Time.Time, 1*time.Minute)

		assert.EqualValues(t, memory, *stat.Memory.UsageBytes)
		assert.WithinDuration(t, time.Now(), stat.Memory.Time.Time, 1*time.Minute)

		assert.EqualValues(t, storage, *stat.VolumeStats[0].UsedBytes)
		assert.WithinDuration(t, time.Now(), stat.VolumeStats[0].Time.Time, 1*time.Minute)

		assert.EqualValues(t, ephStorage, *stat.EphemeralStorage.UsedBytes)
		assert.WithinDuration(t, time.Now(), stat.Memory.Time.Time, 1*time.Minute)

		assert.Equal(t, events.PodStatus, translatePodStatus(podconfigv1.PodStatusRunning))
	})

	ms.EXPECT().SetPodTasks(
		"TestNamespace/TestName",
		gomock.Any(),
	).Do(func(_ string, arr []*store.Task) {
		// Test if the array is empty when no spec tasks are given
		assert.Equal(t, 0, len(arr))
	})

	err := setPodTasks(&ep, &s)
	assert.NoError(t, err)
}
