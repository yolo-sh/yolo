package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func ClosePort(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	portToClose string,
) error {

	closePortErr := cloudService.ClosePort(
		stepper,
		yoloConfig,
		cluster,
		env,
		portToClose,
	)

	// "closePortErr" is not handled first
	// in order to be able to save partial infrastructure
	err := UpdateEnvInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
		env,
	)

	if err != nil {
		return err
	}

	return closePortErr
}
