package git

import (
	"math"
	"regexp"
	"strings"
)

func NewRepoManager(repoPath string, gitExecutor Executor) RepoManager {
	return &repoManager{
		repoPath: repoPath,
		executor: gitExecutor,
	}
}

type RepoManager interface {
	Checkout(branch string) error
	Fetch() error
	FetchAll() error
	RemoteBranches() ([]string, error)
}

type repoManager struct {
	repoPath string
	executor Executor
}

func (repo *repoManager) Checkout(branch string) error {
	return repo.run("checkout", branch)
}

func (repo *repoManager) Fetch() error {
	return repo.run("fetch")
}

func (repo *repoManager) FetchAll() error {
	return repo.run("fetch", "--all")
}

func (repo *repoManager) RemoteBranches() ([]string, error) {
	output, err := repo.output("remote", "-v")
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`(^.+?)\s`)

	var branches []string

	for i, s := range strings.Split(string(output), "\n") {
		if math.Mod(float64(i), 2) != 0 {
			continue
		}
		if s == "" {
			continue
		}
		branches = append(branches, strings.TrimSpace(reg.FindString(s)))
	}
	return branches, nil
}

func (repo *repoManager) run(args ...string) error {
	return repo.executor.RunWithWorkDir(repo.repoPath, args...)
}

func (repo *repoManager) output(args ...string) ([]byte, error) {
	output, err := repo.executor.OutputWithWorkDir(repo.repoPath, args...)
	return output, err
}
