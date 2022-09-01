package actions

import (
	"github.com/yolo-sh/yolo/entities"
	"github.com/yolo-sh/yolo/stepper"
)

func InstallYolo(
	stepper stepper.Stepper,
	cloudService entities.CloudService,
	yoloConfig *entities.Config,
) error {

	err := cloudService.CreateYoloConfigStorage(stepper)

	if err != nil {
		return err
	}

	return cloudService.SaveYoloConfig(
		stepper,
		yoloConfig,
	)
}
