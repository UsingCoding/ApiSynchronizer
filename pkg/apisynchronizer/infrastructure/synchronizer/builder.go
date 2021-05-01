package synchronizer

import (
	"os"
	"path"
)

func NewApiFilePathBuilder(apisRepoPath string) ApiFilePathBuilder {
	return &filePathBuilder{apisRepoPath: apisRepoPath}
}

type ApiFilePathBuilder interface {
	BuildForApiRepo(service Service) string
}

type filePathBuilder struct {
	apisRepoPath string
}

func (builder *filePathBuilder) BuildForApiRepo(service Service) string {
	const apiFileExtension = ".proto"
	return path.Join(builder.apisRepoPath, string(service), string(service)+apiFileExtension)
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
