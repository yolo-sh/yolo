package entities

import "github.com/yolo-sh/yolo/stepper"

type CloudService interface {
	CreateYoloConfigStorage(stepper.Stepper) error
	RemoveYoloConfigStorage(stepper.Stepper) error

	LookupYoloConfig(stepper.Stepper) (*Config, error)
	SaveYoloConfig(stepper.Stepper, *Config) error

	CreateCluster(stepper.Stepper, *Config, *Cluster) error
	RemoveCluster(stepper.Stepper, *Config, *Cluster) error

	CheckInstanceTypeValidity(stepper.Stepper, string) error

	CreateEnv(stepper.Stepper, *Config, *Cluster, *Env) error
	RemoveEnv(stepper.Stepper, *Config, *Cluster, *Env) error

	OpenPort(stepper.Stepper, *Config, *Cluster, *Env, string) error
	ClosePort(stepper.Stepper, *Config, *Cluster, *Env, string) error
}

type CloudServiceBuilder interface {
	Build() (CloudService, error)
}
