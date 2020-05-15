package pod

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/cache"

	v1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"
	mockcache "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/mock_cache_store"
)

const testLabel = "testlabel"

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	csm := mockcache.NewMockStore(ctrl)
	var cm cache.Store = csm

	ep := v1.EmulatedPod{}

	csm.EXPECT().List().Return([]interface{}{
		ep,
	})

	informer := NewInformer(&cm)

	res := informer.List()
	assert.Equal(t, 1, len(res))
	assert.Equal(t, []v1.EmulatedPod{ep}, res)
}

func TestFindExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	csm := mockcache.NewMockStore(ctrl)
	var cm cache.Store = csm

	label := testLabel

	ep := v1.EmulatedPod{}

	csm.EXPECT().GetByKey(label).Return(&ep, true, nil)

	informer := NewInformer(&cm)

	res, found, err := informer.Find(label)
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, &ep, res)
}

func TestFindError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	csm := mockcache.NewMockStore(ctrl)
	var cm cache.Store = csm

	label := testLabel

	se := "TestError"

	csm.EXPECT().GetByKey(label).Return(nil, false, errors.New(se))

	informer := NewInformer(&cm)

	res, found, err := informer.Find(label)
	assert.EqualError(t, err, se)
	assert.False(t, found)
	assert.Nil(t, res)
}

func TestFindNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	csm := mockcache.NewMockStore(ctrl)
	var cm cache.Store = csm

	label := testLabel

	csm.EXPECT().GetByKey(label).Return(nil, false, nil)

	informer := NewInformer(&cm)

	res, found, err := informer.Find(label)
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, res)
}

func TestFindNotCastable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	csm := mockcache.NewMockStore(ctrl)
	var cm cache.Store = csm

	label := testLabel

	type TestStruct struct{}

	csm.EXPECT().GetByKey(label).Return(TestStruct{}, true, nil)

	informer := NewInformer(&cm)

	res, found, err := informer.Find(label)
	assert.Error(t, err)
	assert.False(t, found)
	assert.Nil(t, res)
}
