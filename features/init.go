package features

import (
	"errors"
	"fmt"

	"github.com/yolo-sh/yolo/actions"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

type InitInput struct {
	InstanceType       string
	ResolvedRepository entities.ResolvedEnvRepository
}

type InitOutput struct {
	Error   error
	Content *InitOutputContent
	Stepper stepper.Stepper
}

type InitOutputContent struct {
	CloudService    entities.CloudService
	YoloConfig      *entities.Config
	Cluster         *entities.Cluster
	Env             *entities.Env
	EnvCreated      bool
	SetEnvAsCreated func() error
}

type InitOutputHandler interface {
	HandleOutput(InitOutput) error
}

type InitFeature struct {
	stepper             stepper.Stepper
	outputHandler       InitOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewInitFeature(
	stepper stepper.Stepper,
	outputHandler InitOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) InitFeature {

	return InitFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (i InitFeature) Execute(input InitInput) error {
	handleError := func(err error) error {
		i.outputHandler.HandleOutput(InitOutput{
			Stepper: i.stepper,
			Error:   err,
		})

		return err
	}

	envName := entities.BuildEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	step := fmt.Sprintf("Initializing an environment for \"%s\"", envName)
	i.stepper.StartTemporaryStep(step)

	cloudService, err := i.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	err = cloudService.CheckInstanceTypeValidity(
		i.stepper,
		input.InstanceType,
	)

	if err != nil {
		return handleError(err)
	}

	yoloConfig, err := cloudService.LookupYoloConfig(
		i.stepper,
	)

	if err != nil && !errors.Is(err, entities.ErrYoloNotInstalled) {
		return handleError(err)
	}

	if yoloConfig == nil { // Yolo not installed

		i.stepper.StartTemporaryStep("Installing Yolo")

		yoloConfig = entities.NewConfig()

		err = actions.InstallYolo(
			i.stepper,
			cloudService,
			yoloConfig,
		)

		if err != nil {
			return handleError(err)
		}
	}

	clusterName := entities.DefaultClusterName
	cluster, err := yoloConfig.GetCluster(clusterName)

	if err != nil && !errors.As(err, &entities.ErrClusterNotExists{}) {
		return handleError(err)
	}

	if cluster == nil || cluster.Status == entities.ClusterStatusCreating {

		/* Cluster not exists or still
		in creating state after error */

		i.stepper.StartTemporaryStep("Creating default cluster")

		if cluster == nil {
			// Multiple clusters are not implemented for now
			isDefaultCluster := true

			cluster = entities.NewCluster(
				clusterName,
				input.InstanceType,
				isDefaultCluster,
			)
		}

		err = actions.CreateCluser(
			i.stepper,
			cloudService,
			yoloConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	env, err := yoloConfig.GetEnv(
		cluster.Name,
		envName,
	)

	if err != nil && !errors.As(err, &entities.ErrEnvNotExists{}) {
		return handleError(err)
	}

	if env != nil && env.Status == entities.EnvStatusRemoving {
		return handleError(entities.ErrInitRemovingEnv{
			EnvName: env.Name,
		})
	}

	envCreated := false

	if env == nil || env.Status == entities.EnvStatusCreating {

		/* Env not exists or still
		in creating state after error */

		if env == nil {
			env = entities.NewEnv(
				envName,
				input.InstanceType,
				input.ResolvedRepository,
			)
		}

		err = actions.CreateEnv(
			i.stepper,
			cloudService,
			yoloConfig,
			cluster,
			env,
		)

		if err != nil {
			return handleError(err)
		}

		envCreated = true
	}

	// Current step is the last ended infrastructure step.
	// Better UX if we reset to main step here given that
	// the next steps (in GRPC agent) may take some time to start.
	i.stepper.StartTemporaryStep(step)

	setEnvAsCreated := func() error {
		env.Status = entities.EnvStatusCreated

		return actions.UpdateEnvInConfig(
			i.stepper,
			cloudService,
			yoloConfig,
			cluster,
			env,
		)
	}

	return i.outputHandler.HandleOutput(InitOutput{
		Stepper: i.stepper,
		Content: &InitOutputContent{
			CloudService:    cloudService,
			YoloConfig:      yoloConfig,
			Cluster:         cluster,
			Env:             env,
			EnvCreated:      envCreated,
			SetEnvAsCreated: setEnvAsCreated,
		},
	})
}
