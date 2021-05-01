package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseConfig() (*config, error) {
	c := &config{
		ApisRepoUrl:   "git@github.com:CuriosityMusicStreaming/ApiGateway.git",
		RepoCachePath: "/tmp/apisynchronizer/cache",
	}
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ApisRepoUrl   string `envconfig:"apis_repo_url"`
	RepoCachePath string `envconfig:"repo_cache_path"`
}
