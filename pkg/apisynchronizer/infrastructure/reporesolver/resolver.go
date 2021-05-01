package reporesolver

import (
	git2 "apisynchronizer/pkg/common/infrastructure/git"
	"github.com/pkg/errors"
	"os"
	"path"
)

const (
	gitFolder     = ".git"
	cacheRepoName = "@apis"
)

var (
	errFolderExistsButItNotARepo = errors.New("folder repo exists but its not a git repo")
)

type Resolver interface {
	Path() (string, error)
}

func New(apisRepoUrl, apisRepoCachePath string, gitExecutor git2.Executor) Resolver {
	return &repoResolver{
		apisRepoUrl:       apisRepoUrl,
		apisRepoCachePath: apisRepoCachePath,
		gitExecutor:       gitExecutor,
	}
}

type repoResolver struct {
	apisRepoUrl       string
	apisRepoCachePath string
	gitExecutor       git2.Executor
	repoResolved      bool
}

func (resolver *repoResolver) Path() (string, error) {
	if !resolver.repoResolved {
		err := resolver.resolveRepo()
		if err != nil {
			return "", err
		}
	}

	return path.Join(resolver.apisRepoCachePath, cacheRepoName), nil
}

func (resolver *repoResolver) resolveRepo() error {
	err := repoInCacheExists(path.Join(resolver.apisRepoCachePath, cacheRepoName))
	if err == errFolderExistsButItNotARepo {
		return err
	}
	if err != nil {
		err2 := os.MkdirAll(resolver.apisRepoCachePath, 0755)
		if err2 != nil {
			return errors.Wrap(err2, "failed to create repo cache folder")
		}
		err2 = resolver.gitExecutor.RunWithWorkDir(resolver.apisRepoCachePath, "clone", resolver.apisRepoUrl, cacheRepoName)
		if err2 != nil {
			return errors.Wrap(err2, "failed to clone repo")
		}
	}

	return nil
}

func repoInCacheExists(repoPath string) error {
	stat, err := os.Stat(repoPath)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("repo path is not a folder")
	}

	_, err = os.Stat(path.Join(repoPath, gitFolder))
	if err != nil {
		return errors.Wrap(errFolderExistsButItNotARepo, err.Error())
	}
	return nil
}
