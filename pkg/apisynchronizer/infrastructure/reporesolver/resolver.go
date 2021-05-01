package reporesolver

import (
	embedder "apisynchronizer/data/apisynchronizer"
	"apisynchronizer/pkg/apisynchronizer/infrastructure"
	"apisynchronizer/pkg/common/infrastructure/git"
	"github.com/pkg/errors"
	"io/ioutil"
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

func New(
	apisRepoUrl,
	apisRepoCachePath string,
	gitExecutor git.Executor,
	reporter infrastructure.Reporter,
) Resolver {
	return &repoResolver{
		apisRepoUrl:       apisRepoUrl,
		apisRepoCachePath: apisRepoCachePath,
		gitExecutor:       gitExecutor,
		reporter:          reporter,
	}
}

type repoResolver struct {
	apisRepoUrl       string
	apisRepoCachePath string
	gitExecutor       git.Executor
	reporter          infrastructure.Reporter
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
		resolver.reporter.Info("No local repo, cloning repo ⏳...")
		err2 := os.MkdirAll(resolver.apisRepoCachePath, 0755)
		if err2 != nil {
			return errors.Wrap(err2, "failed to create repo cache folder")
		}
		err2 = resolver.gitExecutor.RunWithWorkDir(resolver.apisRepoCachePath, "clone", resolver.apisRepoUrl, cacheRepoName)
		if err2 != nil {
			return errors.Wrap(err2, "failed to clone repo")
		}

		err2 = ioutil.WriteFile(path.Join(resolver.apisRepoCachePath, cacheRepoName, embedder.WarningReadMeName), embedder.WarningReadMe, 0755)
		if err2 != nil {
			return err2
		}

		resolver.reporter.Info("Remote repo cloned 🔥")
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
