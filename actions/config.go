package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func UpdateClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := yoloConfig.SetCluster(cluster)

	if err != nil {
		return err
	}

	return cloudService.SaveYoloConfig(
		stepper,
		yoloConfig,
	)
}

func RemoveClusterInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	err := yoloConfig.RemoveCluster(cluster.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveYoloConfig(
		stepper,
		yoloConfig,
	)
}

func UpdateEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	err := yoloConfig.SetEnv(cluster.Name, env)

	if err != nil {
		return err
	}

	return cloudService.SaveYoloConfig(
		stepper,
		yoloConfig,
	)
}

func RemoveEnvInConfig(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
	env *entities.Env,
) error {

	err := yoloConfig.RemoveEnv(cluster.Name, env.Name)

	if err != nil {
		return err
	}

	return cloudService.SaveYoloConfig(
		stepper,
		yoloConfig,
	)
}
