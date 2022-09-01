package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	giturls "github.com/whilp/git-urls"
	"github.com/yolo-sh/yolo/entities"
)

func BuildGitHTTPURL(repoOwner, repoName string) entities.EnvRepositoryGitURL {
	return entities.EnvRepositoryGitURL(fmt.Sprintf(
		"https://github.com/%s/%s.git",
		url.PathEscape(repoOwner),
		url.PathEscape(repoName),
	))
}

func BuildGitURL(repoOwner, repoName string) entities.EnvRepositoryGitURL {
	return entities.EnvRepositoryGitURL(fmt.Sprintf(
		"git@github.com:%s/%s.git",
		url.PathEscape(repoOwner),
		url.PathEscape(repoName),
	))
}

type ParsedGitHubRepositoryName struct {
	Owner         string
	ExplicitOwner bool
	Name          string
}

func ParseRepositoryName(
	repositoryName string,
	defaultRepositoryOwner string,
) (*ParsedGitHubRepositoryName, error) {

	errInvalidGitHubURL := errors.New("ErrInvalidGitHubURL")

	// Handle git@github.com:yolo-sh/yolo.git
	repositoryNameAsURL, err := giturls.Parse(repositoryName)

	if err != nil {
		// Handle https://github.com/yolo-sh/yolo.git
		repositoryNameAsURL, err = url.Parse(repositoryName)
	}

	// Not an URL (eg: yolo) or only path (eg: yolo-sh/yolo)
	if err != nil || len(repositoryNameAsURL.Hostname()) == 0 {
		repositoryNameParts := strings.Split(repositoryName, "/")

		if len(repositoryNameParts) > 2 {
			return nil, errInvalidGitHubURL
		}

		if len(repositoryNameParts) == 1 { // yolo
			return &ParsedGitHubRepositoryName{
				ExplicitOwner: false,
				Owner:         defaultRepositoryOwner,
				Name:          repositoryNameParts[0],
			}, nil
		}

		return &ParsedGitHubRepositoryName{ // yolo-sh/yolo
			ExplicitOwner: true,
			Owner:         repositoryNameParts[0],
			Name:          repositoryNameParts[1],
		}, nil
	}

	host := repositoryNameAsURL.Hostname()

	if host != "github.com" {
		return nil, errInvalidGitHubURL
	}

	path := strings.TrimPrefix(repositoryNameAsURL.Path, "/")
	pathComponents := strings.Split(path, "/")

	if len(pathComponents) < 2 {
		return nil, errInvalidGitHubURL
	}

	githubRepositoryOwner := pathComponents[0]
	githubRepositoryName := strings.TrimSuffix(pathComponents[1], ".git")

	return &ParsedGitHubRepositoryName{
		ExplicitOwner: true,
		Owner:         githubRepositoryOwner,
		Name:          githubRepositoryName,
	}, nil
}
