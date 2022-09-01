package features

import (
	"fmt"

	"github.com/yolo-sh/yolo/actions"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

type ClosePortInput struct {
	ResolvedRepository entities.ResolvedEnvRepository
	PortToClose        string
}

type ClosePortOutput struct {
	Error   error
	Content *ClosePortOutputContent
	Stepper stepper.Stepper
}

type ClosePortOutputContent struct {
	Cluster           *entities.Cluster
	Env               *entities.Env
	PortClosed        string
	PortAlreadyClosed bool
}

type ClosePortOutputHandler interface {
	HandleOutput(ClosePortOutput) error
}

type ClosePortFeature struct {
	stepper             stepper.Stepper
	outputHandler       ClosePortOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewClosePortFeature(
	stepper stepper.Stepper,
	outputHandler ClosePortOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) ClosePortFeature {

	return ClosePortFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (o ClosePortFeature) Execute(input ClosePortInput) error {
	handleError := func(err error) error {
		o.outputHandler.HandleOutput(ClosePortOutput{
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
			"Closing port \"%s\"",
			input.PortToClose,
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
		return handleError(entities.ErrClosePortRemovingEnv{
			EnvName: envName,
		})
	}

	if env.Status == entities.EnvStatusCreating {
		return handleError(entities.ErrClosePortCreatingEnv{
			EnvName: envName,
		})
	}

	portAlreadyClosed := !env.OpenedPorts[input.PortToClose]

	if !portAlreadyClosed {
		err = actions.ClosePort(
			o.stepper,
			cloudService,
			yoloConfig,
			cluster,
			env,
			input.PortToClose,
		)

		if err != nil {
			return handleError(err)
		}
	}

	return o.outputHandler.HandleOutput(ClosePortOutput{
		Stepper: o.stepper,
		Content: &ClosePortOutputContent{
			Cluster:           cluster,
			Env:               env,
			PortClosed:        input.PortToClose,
			PortAlreadyClosed: portAlreadyClosed,
		},
	})
}
