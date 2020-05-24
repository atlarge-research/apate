package run

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/run/mock_run"
)

const MyRunType env.RunType = "newRunType"

func setEnv() {
	once := sync.Once{}
	once.Do(func() {
		planeEnv := env.ControlPlaneEnv()
		planeEnv.ApateletRunType = MyRunType
		env.SetEnv(planeEnv)
	})
}

func TestRegisterRunner(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mar := mock_run.NewMockApateletRunner(ctrl)
	var r ApateletRunner = mar

	registry := New()
	registry.RegisterRunner(MyRunType, &r, 1, 2)

	ctx := context.TODO()
	environment := env.ApateletEnvironment{}

	mar.EXPECT().SpawnApatelets(ctx, 20, environment, 1, 2).Return(nil)

	err := registry.Run(ctx, 20, environment)
	assert.NoError(t, err)
}

func TestRegisterRunnerUnknownType(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registry := New()
	err := registry.Run(context.TODO(), 20, env.ApateletEnvironment{})
	assert.Error(t, err)
}

func TestRegisterRunnerReturnsError(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mar := mock_run.NewMockApateletRunner(ctrl)
	var r ApateletRunner = mar

	registry := New()
	registry.RegisterRunner(MyRunType, &r, 1, 2)

	ctx := context.TODO()
	environment := env.ApateletEnvironment{}

	mar.EXPECT().SpawnApatelets(ctx, 20, environment, 1, 2).Return(errors.New("oops"))

	err := registry.Run(ctx, 20, environment)
	assert.Error(t, err)
}
