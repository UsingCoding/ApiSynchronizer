package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	apiRepoUrl = "UNKNOWN"
)

func parseConfig() (*config, error) {
	c := &config{
		ApisRepoUrl:   apiRepoUrl,
		RepoCachePath: "/tmp/apisynchronizer/cache",
		ApisFolder:    "api",
	}
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

func parseRuntimeConfig(ctx *cli.Context) (*runtimeConfig, error) {
	c := new(runtimeConfig)

	c.outputPath = ctx.String("output")
	c.configPath = ctx.String("file")

	return c, nil
}

type config struct {
	ApisRepoUrl   string `envconfig:"apis_repo_url"`
	RepoCachePath string `envconfig:"repo_cache_path"`
	ApisFolder    string `envconfig:"apis_folder"`
}

type runtimeConfig struct {
	outputPath string
	configPath string
}
