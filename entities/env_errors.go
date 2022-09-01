package entities

type ErrEnvNotExists struct {
	ClusterName string
	EnvName     string
}

func (ErrEnvNotExists) Error() string {
	return "ErrEnvNotExists"
}

type ErrInvalidPort struct {
	InvalidPort string
}

func (ErrInvalidPort) Error() string {
	return "ErrInvalidPort"
}

type ErrReservedPort struct {
	ReservedPort string
}

func (ErrReservedPort) Error() string {
	return "ErrReservedPort"
}

type ErrInitRemovingEnv struct {
	EnvName string
}

func (ErrInitRemovingEnv) Error() string {
	return "ErrInitRemovingEnv"
}

type ErrEditRemovingEnv struct {
	EnvName string
}

func (ErrEditRemovingEnv) Error() string {
	return "ErrEditRemovingEnv"
}

type ErrEditCreatingEnv struct {
	EnvName string
}

func (ErrEditCreatingEnv) Error() string {
	return "ErrEditCreatingEnv"
}

type ErrOpenPortRemovingEnv struct {
	EnvName string
}

func (ErrOpenPortRemovingEnv) Error() string {
	return "ErrOpenPortRemovingEnv"
}

type ErrOpenPortCreatingEnv struct {
	EnvName string
}

func (ErrOpenPortCreatingEnv) Error() string {
	return "ErrOpenPortCreatingEnv"
}

type ErrClosePortRemovingEnv struct {
	EnvName string
}

func (ErrClosePortRemovingEnv) Error() string {
	return "ErrClosePortRemovingEnv"
}

type ErrClosePortCreatingEnv struct {
	EnvName string
}

func (ErrClosePortCreatingEnv) Error() string {
	return "ErrClosePortCreatingEnv"
}
