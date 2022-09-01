package features

import (
	"errors"

	"github.com/yolo-sh/yolo/actions"
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

type UninstallInput struct {
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutput struct {
	Error   error
	Content *UninstallOutputContent
	Stepper stepper.Stepper
}

type UninstallOutputContent struct {
	YoloAlreadyUninstalled    bool
	SuccessMessage            string
	AlreadyUninstalledMessage string
}

type UninstallOutputHandler interface {
	HandleOutput(UninstallOutput) error
}

type UninstallFeature struct {
	stepper             stepper.Stepper
	outputHandler       UninstallOutputHandler
	cloudServiceBuilder entities.CloudServiceBuilder
}

func NewUninstallFeature(
	stepper stepper.Stepper,
	outputHandler UninstallOutputHandler,
	cloudServiceBuilder entities.CloudServiceBuilder,
) UninstallFeature {

	return UninstallFeature{
		stepper:             stepper,
		outputHandler:       outputHandler,
		cloudServiceBuilder: cloudServiceBuilder,
	}
}

func (u UninstallFeature) Execute(input UninstallInput) error {
	handleError := func(err error) error {
		u.outputHandler.HandleOutput(UninstallOutput{
			Stepper: u.stepper,
			Error:   err,
		})

		return err
	}

	u.stepper.StartTemporaryStep("Uninstalling Yolo")

	cloudService, err := u.cloudServiceBuilder.Build()

	if err != nil {
		return handleError(err)
	}

	yoloConfig, err := cloudService.LookupYoloConfig(
		u.stepper,
	)

	if err != nil {
		if errors.Is(err, entities.ErrYoloNotInstalled) {
			return u.outputHandler.HandleOutput(UninstallOutput{
				Stepper: u.stepper,
				Content: &UninstallOutputContent{
					YoloAlreadyUninstalled:    true,
					SuccessMessage:            input.SuccessMessage,
					AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
				},
			})
		}

		return handleError(err)
	}

	clusterName := entities.DefaultClusterName
	cluster, err := yoloConfig.GetCluster(clusterName)

	// In case of error the yolo config storage
	// could be created but without cluster
	if err != nil && !errors.As(err, &entities.ErrClusterNotExists{}) {
		return handleError(err)
	}

	if cluster != nil {
		nbOfEnvsInCluster, err := yoloConfig.CountEnvsInCluster(clusterName)

		if err != nil {
			return handleError(err)
		}

		if nbOfEnvsInCluster > 0 {
			return handleError(entities.ErrUninstallExistingEnvs)
		}

		err = actions.RemoveCluster(
			u.stepper,
			cloudService,
			yoloConfig,
			cluster,
		)

		if err != nil {
			return handleError(err)
		}
	}

	err = cloudService.RemoveYoloConfigStorage(
		u.stepper,
	)

	if err != nil {
		return handleError(err)
	}

	return u.outputHandler.HandleOutput(UninstallOutput{
		Stepper: u.stepper,
		Content: &UninstallOutputContent{
			YoloAlreadyUninstalled:    false,
			SuccessMessage:            input.SuccessMessage,
			AlreadyUninstalledMessage: input.AlreadyUninstalledMessage,
		},
	})
}
