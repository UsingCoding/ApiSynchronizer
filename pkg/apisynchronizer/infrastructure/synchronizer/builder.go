package synchronizer

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func NewApiFileFinder(apisRepoPath string, apisFolder string) ApiFileFinder {
	return &apiFileFinder{
		apisRepoPath: apisRepoPath,
		apisFolder:   apisFolder,
	}
}

type ApiFileFinder interface {
	FindApiForService(service Service) (string, error)
}

type apiFileFinder struct {
	apisRepoPath string
	apisFolder   string
}

func (finder *apiFileFinder) FindApiForService(service Service) (string, error) {
	apisFolderPath := path.Join(finder.apisRepoPath, finder.apisFolder)
	entries, err := ioutil.ReadDir(apisFolderPath)
	if err != nil {
		return "", errors.Wrap(err, "no apis folder found in repo")
	}

	for _, entry := range entries {
		fileNameWithoutExt := strings.Trim(entry.Name(), path.Ext(entry.Name()))
		if fileNameWithoutExt == string(service) {
			return path.Join(apisFolderPath, entry.Name()), nil
		}
	}

	return "", errors.New(fmt.Sprintf("api file for service %s not found", service))
}

func NewOutputStructureBuilder(outputPath string) OutputStructureBuilder {
	return &outputStructureBuilder{outputPath: outputPath}
}

type OutputStructureBuilder interface {
	// Build builds folder structure and returns path to write file there
	Build(service Service) (string, error)
}

type outputStructureBuilder struct {
	outputPath string
}

func (o *outputStructureBuilder) Build(service Service) (string, error) {
	apiPath := path.Join(o.outputPath, string(service))
	err := os.MkdirAll(apiPath, 0755)
	return apiPath, err
}
