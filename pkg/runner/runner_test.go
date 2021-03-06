package runner

import (
	"context"
	"errors"
	"testing"

	"github.com/atlarge-research/apate/pkg/runner/mock_runner"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/atlarge-research/apate/pkg/env"
)

const MyRunType env.RunType = "newRunType"

func setEnv() {
	planeEnv := env.ControlPlaneEnv()
	planeEnv.ApateletRunType = MyRunType
	env.SetEnv(planeEnv)
}

func TestRegisterRunner(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mar := mock_runner.NewMockApateletRunner(ctrl)
	var r ApateletRunner = mar

	registry := New()
	registry.RegisterRunner(MyRunType, &r)

	ctx := context.Background()
	environment := env.ApateletEnvironment{}

	mar.EXPECT().SpawnApatelets(ctx, 20, environment).Return(nil)

	err := registry.Run(ctx, 20, environment)
	assert.NoError(t, err)
}

func TestRegisterRunnerUnknownType(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	registry := New()
	err := registry.Run(context.Background(), 20, env.ApateletEnvironment{})
	assert.Error(t, err)
}

func TestRegisterRunnerReturnsError(t *testing.T) {
	setEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mar := mock_runner.NewMockApateletRunner(ctrl)
	var r ApateletRunner = mar

	registry := New()
	registry.RegisterRunner(MyRunType, &r)

	ctx := context.Background()
	environment := env.ApateletEnvironment{}

	mar.EXPECT().SpawnApatelets(ctx, 20, environment).Return(errors.New("oops"))

	err := registry.Run(ctx, 20, environment)
	assert.Error(t, err)
}
