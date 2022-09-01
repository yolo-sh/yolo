package entities

type HookRunner interface {
	Run(
		cloudService CloudService,
		config *Config,
		cluster *Cluster,
		env *Env,
	) error
}
