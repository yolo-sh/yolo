package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func CreateCluser(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	createClusterErr := cloudService.CreateCluster(
		stepper,
		yoloConfig,
		cluster,
	)

	// "createCLusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if createClusterErr != nil {
		return createClusterErr
	}

	cluster.Status = entities.ClusterStatusCreated
	return UpdateClusterInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
	)
}
