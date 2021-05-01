package synchronizer

import (
	"apisynchronizer/pkg/apisynchronizer/infrastructure"
	"apisynchronizer/pkg/common/infrastructure/git"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
)

type (
	Service  string
	Revision string
)

func New(
	repoManager git.RepoManager,
	pathBuilder ApiFilePathBuilder,
	structureBuilder OutputStructureBuilder,
	reporter infrastructure.Reporter,
) *Synchronizer {
	return &Synchronizer{
		repoManager:      repoManager,
		pathBuilder:      pathBuilder,
		structureBuilder: structureBuilder,
		reporter:         reporter,
	}
}

type Synchronizer struct {
	repoManager      git.RepoManager
	pathBuilder      ApiFilePathBuilder
	structureBuilder OutputStructureBuilder
	reporter         infrastructure.Reporter
}

type SynchronizeParams struct {
	ApiDeclaration map[Service]Revision
}

func (s *Synchronizer) Synchronize(params SynchronizeParams) error {
	remoteBranch, err := s.fetchRemoteBranch()
	if err != nil {
		return err
	}

	err = s.repoManager.FetchAll()
	if err != nil {
		return err
	}

	serviceApiFile := map[Service]string{}

	for service, revision := range params.ApiDeclaration {
		err2 := s.repoManager.Checkout(fmt.Sprintf("%s/%s", remoteBranch, revision))
		if err2 != nil {
			return err2
		}

		apiFileInApiRepo := s.pathBuilder.BuildForApiRepo(service)
		if err2 = fileExists(apiFileInApiRepo); err2 != nil {
			return errors.Wrap(err2, fmt.Sprintf("service api file for %s doesnt exists in api repo", service))
		}

		serviceApiFile[service] = apiFileInApiRepo

		s.reporter.Info(fmt.Sprintf("Api file for %s resolved", service))
	}

	s.reporter.Info("All api files resolved 👹")

	for service, apiFilePath := range serviceApiFile {
		outputPath, err2 := s.structureBuilder.Build(service)
		if err2 != nil {
			return err2
		}

		fileName := path.Base(apiFilePath)

		newApiFilePath := path.Join(outputPath, fileName)
		err2 = copyFile(newApiFilePath, apiFilePath)
		if err2 != nil {
			return err2
		}

		s.reporter.Info(fmt.Sprintf("Api file for %s synced to %s", service, newApiFilePath))
	}

	return nil
}

func (s *Synchronizer) fetchRemoteBranch() (string, error) {
	branches, err := s.repoManager.RemoteBranches()
	if err != nil {
		return "", err
	}

	if len(branches) != 1 {
		return "", errors.New(fmt.Sprintf("to many remote branches: %s", branches))
	}

	return branches[0], nil
}

func fileExists(path string) error {
	_, err := os.Stat(path)
	return err
}

func copyFile(dst, source string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0644)
	return err
}
