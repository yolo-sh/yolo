package entities

import "errors"

var (
	ErrYoloNotInstalled      = errors.New("ErrYoloNotInstalled")
	ErrUninstallExistingEnvs = errors.New("ErrUninstallExistingEnvs")
)
