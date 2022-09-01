package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func CreateEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	createEnvErr := cloudService.CreateEnv(
		stepper,
		yoloConfig,
		cluster,
		env,
	)

	// "createEnvErr" is not handled first
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

	return createEnvErr
}
