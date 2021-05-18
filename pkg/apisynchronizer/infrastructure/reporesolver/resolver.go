package reporesolver

import (
	"apisynchronizer/pkg/common/infrastructure/git"
	"apisynchronizer/pkg/common/infrastructure/reporter"
	"github.com/pkg/errors"
	"os"
	"path"
)

const (
	gitFolder = ".git"
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
	reporter reporter.Reporter,
) Resolver {
	return &repoResolver{
		apisRepoUrl:  apisRepoUrl,
		apisRepoPath: apisRepoCachePath,
		gitExecutor:  gitExecutor,
		reporter:     reporter,
	}
}

type repoResolver struct {
	apisRepoUrl  string
	apisRepoPath string
	gitExecutor  git.Executor
	reporter     reporter.Reporter
	repoResolved bool
}

func (resolver *repoResolver) Path() (string, error) {
	if !resolver.repoResolved {
		err := resolver.resolveRepo()
		if err != nil {
			return "", err
		}
	}

	return resolver.apisRepoPath, nil
}

func (resolver *repoResolver) resolveRepo() error {
	err := repoInCacheExists(resolver.apisRepoPath)
	if err == errFolderExistsButItNotARepo {
		return err
	}
	if err != nil {
		resolver.reporter.Info("No local repo, cloning repo ⏳...")
		err2 := os.MkdirAll(resolver.apisRepoPath, 0755)
		if err2 != nil {
			return errors.Wrap(err2, "failed to create repo cache folder")
		}
		repoFolder, repoName := path.Split(resolver.apisRepoPath)
		err2 = resolver.gitExecutor.RunWithWorkDir(repoFolder, "clone", resolver.apisRepoUrl, repoName)
		if err2 != nil {
			return errors.Wrap(err2, "failed to clone repo")
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
