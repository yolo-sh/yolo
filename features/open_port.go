package features

import (
	"fmt"

	"github.com/yolo-sh/yolo/actions"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

type OpenPortInput struct {
	ResolvedRepository entities.ResolvedEnvRepository
	PortToOpen         string
}

type OpenPortOutput struct {
	Error   error
	Content *OpenPortOutputContent
	Stepper stepper.Stepper
}

type OpenPortOutputContent struct {
	Cluster           *entities.Cluster
	Env               *entities.Env
	PortOpened        string
	PortAlreadyOpened bool
}

type OpenPortOutputHandler interface {
	HandleOutput(OpenPortOutput) error
}

type OpenPortFeature struct {
	stepper             stepper.Stepper
	outputHandler       OpenPortOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewOpenPortFeature(
	stepper stepper.Stepper,
	outputHandler OpenPortOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) OpenPortFeature {

	return OpenPortFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (o OpenPortFeature) Execute(input OpenPortInput) error {
	handleError := func(err error) error {
		o.outputHandler.HandleOutput(OpenPortOutput{
			Stepper: o.stepper,
			Error:   err,
		})

		return err
	}

	envName := entities.BuildEnvNameFromResolvedRepo(
		input.ResolvedRepository,
	)

	o.stepper.StartTemporaryStep(
		fmt.Sprintf(
			"Opening port \"%s\"",
			input.PortToOpen,
		),
	)

	cloudService, err := o.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	yoloConfig, err := cloudService.LookupYoloConfig(
		o.stepper,
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

	if env.Status == entities.EnvStatusRemoving {
		return handleError(entities.ErrOpenPortRemovingEnv{
			EnvName: envName,
		})
	}

	if env.Status == entities.EnvStatusCreating {
		return handleError(entities.ErrOpenPortCreatingEnv{
			EnvName: envName,
		})
	}

	portAlreadyOpened := env.OpenedPorts[input.PortToOpen]

	if !portAlreadyOpened {
		err = actions.OpenPort(
			o.stepper,
			cloudService,
			yoloConfig,
			cluster,
			env,
			input.PortToOpen,
		)

		if err != nil {
			return handleError(err)
		}
	}

	return o.outputHandler.HandleOutput(OpenPortOutput{
		Stepper: o.stepper,
		Content: &OpenPortOutputContent{
			Cluster:           cluster,
			Env:               env,
			PortOpened:        input.PortToOpen,
			PortAlreadyOpened: portAlreadyOpened,
		},
	})
}
