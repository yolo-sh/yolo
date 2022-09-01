package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func RemoveCluster(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
	cluster *entities.Cluster,
) error {

	cluster.Status = entities.ClusterStatusRemoving
	err := UpdateClusterInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	removeClusterErr := cloudService.RemoveCluster(
		stepper,
		yoloConfig,
		cluster,
	)

	// "removeClusterErr" is not handled first
	// in order to be able to save partial infrastructure
	err = UpdateClusterInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
	)

	if err != nil {
		return err
	}

	if removeClusterErr != nil {
		return removeClusterErr
	}

	return RemoveClusterInConfig(
		stepper,
		cloudService,
		yoloConfig,
		cluster,
	)
}
