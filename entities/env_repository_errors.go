package entities

type ErrEnvRepositoryNotFound struct {
	RepoOwner string
	RepoName  string
}

func (ErrEnvRepositoryNotFound) Error() string {
	return "ErrEnvRepositoryNotFound"
}
