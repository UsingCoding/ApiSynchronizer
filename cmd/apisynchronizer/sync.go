package main

import (
	"apisynchronizer/pkg/apisynchronizer/infrastructure/reporesolver"
	"apisynchronizer/pkg/apisynchronizer/infrastructure/synchronizer"
	"apisynchronizer/pkg/common/infrastructure/git"
	"apisynchronizer/pkg/common/infrastructure/reporter"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func executeResolve(ctx *cli.Context) error {
	config, err := parseConfig()
	if err != nil {
		return err
	}

	runtimeConfig, err := parseRuntimeConfig(ctx)
	if err != nil {
		return err
	}

	apiDeclarations, err := parseApiDeclarations(runtimeConfig.configPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	gitExecutor, err := git.NewExecutor()
	if err != nil {
		return err
	}

	infoReporter := initReporter(ctx)

	repoResolver := reporesolver.New(
		config.ApisRepoUrl,
		config.LocalRepoPath,
		gitExecutor,
		infoReporter,
	)

	apisRepoPath, err := repoResolver.Path()
	if err != nil {
		return err
	}

	repoManager := git.NewRepoManager(apisRepoPath, gitExecutor)

	apiSynchronizer := initSynchronizer(apisRepoPath, config.ApisFolder, runtimeConfig.outputPath, repoManager, infoReporter)

	return apiSynchronizer.Synchronize(synchronizer.SynchronizeParams{
		ApiDeclaration: apiDeclarations,
		ForceRemote:    runtimeConfig.forceRemote,
	})
}

func initSynchronizer(
	apisRepoPath,
	apisFolder,
	outputPath string,
	manager git.RepoManager,
	reporter reporter.Reporter,
) *synchronizer.Synchronizer {
	return synchronizer.New(
		manager,
		synchronizer.NewApiFileFinder(apisRepoPath, apisFolder),
		synchronizer.NewOutputStructureBuilder(outputPath),
		reporter,
	)
}

func convertMap(input map[string]string) map[synchronizer.Service]synchronizer.Revision {
	output := map[synchronizer.Service]synchronizer.Revision{}
	for service, revision := range input {
		output[synchronizer.Service(service)] = synchronizer.Revision(revision)
	}
	return output
}

func parseApiDeclarations(path string) (map[synchronizer.Service]synchronizer.Revision, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	type declarationsConfig struct {
		Apis map[string]string `yaml:"api"`
	}

	c := new(declarationsConfig)

	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return nil, err
	}

	return convertMap(c.Apis), nil
}
