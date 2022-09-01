package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func RemoveEnv(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
	preRemoveHook entities.HookRunner,
) error {

	env.Status = entities.EnvStatusRemoving
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

	removeEnvErr := cloudService.RemoveEnv(
		stepper,
		yoloConfig,
		cluster,
		env,
	)

	// "removeEnvErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateEnvInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
		env,
	)

	if err != nil {
		return err
	}

	if removeEnvErr != nil {
		return removeEnvErr
	}

	if preRemoveHook != nil {
		err = preRemoveHook.Run(
			cloudService,
			yoloConfig,
			cluster,
			env,
		)

		if err != nil {
			return err
		}
	}

	return RemoveEnvInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
		env,
	)
}
