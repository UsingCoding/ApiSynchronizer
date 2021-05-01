package synchronizer

import (
	"apisynchronizer/pkg/common/infrastructure/git"
	"apisynchronizer/pkg/common/infrastructure/reporter"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path"
)

type (
	Service  string
	Revision string
)

func New(
	repoManager git.RepoManager,
	pathBuilder ApiFileFinder,
	structureBuilder OutputStructureBuilder,
	reporter reporter.Reporter,
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
	pathBuilder      ApiFileFinder
	structureBuilder OutputStructureBuilder
	reporter         reporter.Reporter
}

type SynchronizeParams struct {
	ApiDeclaration map[Service]Revision
}

func (s *Synchronizer) Synchronize(params SynchronizeParams) error {
	remoteBranch, err := s.fetchRemoteBranch()
	if err != nil {
		return err
	}

	s.reporter.Info("Fetching repo ⏳...")

	err = s.repoManager.FetchAll()
	if err != nil {
		return err
	}

	serviceApiFile := map[Service]struct {
		fileName string
		file     []byte
	}{}

	for service, revision := range params.ApiDeclaration {
		err2 := s.repoManager.Checkout(fmt.Sprintf("%s/%s", remoteBranch, revision))
		if err2 != nil {
			return err2
		}

		apiFileInApiRepo, err2 := s.pathBuilder.FindApiForService(service)
		if err2 != nil {
			return err2
		}

		input, err := ioutil.ReadFile(apiFileInApiRepo)
		if err != nil {
			return err
		}

		serviceApiFile[service] = struct {
			fileName string
			file     []byte
		}{fileName: path.Base(apiFileInApiRepo), file: input}

		s.reporter.Info(fmt.Sprintf("Api file for %s resolved", service))
	}

	s.reporter.Info("All api files resolved 👹")

	for service, fileInfo := range serviceApiFile {
		outputPath, err2 := s.structureBuilder.Build(service)
		if err2 != nil {
			return err2
		}

		newApiFilePath := path.Join(outputPath, fileInfo.fileName)

		err2 = ioutil.WriteFile(newApiFilePath, fileInfo.file, 0755)
		if err2 != nil {
			return err2
		}

		s.reporter.Info(fmt.Sprintf("Api file for %s synced to %s", service, newApiFilePath))
	}

	s.reporter.Info("All api files synchronized ✅")

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
