package features

import (
	"fmt"

	"github.com/yolo-sh/yolo/actions"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

type RemoveInput struct {
	ResolvedRepository entities.ResolvedEnvRepository
	PreRemoveHook      entities.HookRunner
	ForceRemove        bool
	ConfirmRemove      func() (bool, error)
}

type RemoveOutput struct {
	Error   error
	Content *RemoveOutputContent
	Stepper stepper.Stepper
}

type RemoveOutputContent struct {
	Cluster *entities.Cluster
	Env     *entities.Env
}

type RemoveOutputHandler interface {
	HandleOutput(RemoveOutput) error
}

type RemoveFeature struct {
	stepper             stepper.Stepper
	outputHandler       RemoveOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewRemoveFeature(
	stepper stepper.Stepper,
	outputHandler RemoveOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) RemoveFeature {

	return RemoveFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (r RemoveFeature) Execute(input RemoveInput) error {
	handleError := func(err error) error {
		r.outputHandler.HandleOutput(RemoveOutput{
			Stepper: r.stepper,
			Error:   err,
		})

		return err
	}

	envName := entities.BuildEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	step := fmt.Sprintf("Removing the environment for \"%s\"", envName)
	r.stepper.StartTemporaryStep(step)

	cloudService, err := r.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	yoloConfig, err := cloudService.LookupYoloConfig(
		r.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := yoloConfig.GetCluster(clusterName)

	if err != nil {
		return handleError(err)
	}

	env, err := yoloConfig.GetEnv(cluster.Name, envName)

	if err != nil {
		return handleError(err)
	}

	if !input.ForceRemove && input.ConfirmRemove != nil {
		r.stepper.StopCurrentStep()

		confirmed, err := input.ConfirmRemove()

		if err != nil {
			return handleError(err)
		}

		if !confirmed {
			return nil
		}

		r.stepper.StartTemporaryStep(step)
	}

	err = actions.RemoveEnv(
		r.stepper,
		cloudService,
		yoloConfig,
		cluster,
		env,
		input.PreRemoveHook,
	)

	if err != nil {
		return handleError(err)
	}

	return r.outputHandler.HandleOutput(RemoveOutput{
		Stepper: r.stepper,
		Content: &RemoveOutputContent{
			Cluster: cluster,
			Env:     env,
		},
	})
}
