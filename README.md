# Yolo

This repository contains all the entities, actions and features used by all the other packages (including the [CLI](https://github.com/yolo-sh/cli)). 

*For the clean architecture aficionados, we are in the innermost circle: the entities one.*

## Table of contents
- [Requirements](#requirements)
- [Usage](#usage)
- [License](#license)

## Requirements

- `go >= 1.18` (this module makes use of generics)

## Usage

**This repository is not meant to be used standalone (you could see that there is no `main.go` file). It is only meant to be imported by other packages**.

As an example, all the cloud providers added to the [CLI](https://github.com/yolo-sh/cli) need to conform to the `CloudService` interface:

```go
// entities/cloud_service.go
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
```

## License

Yolo is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
