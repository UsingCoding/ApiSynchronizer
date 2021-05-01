package main

import (
	"apisynchronizer/pkg/common/infrastructure/reporter"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	stdlog "log"
	"os"
	"time"
)

var (
	appID   = "UNKNOWN"
	version = "UNKNOWN"
)

func main() {
	err := runApp(os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(args []string) error {
	app := &cli.App{
		Name:                 appID,
		Version:              version,
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:      "sync",
				UsageText: "sync -o <output>",
				Usage:     "Syncs api-files for service",
				Action:    executeResolve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Required: true,
					},
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Value:       "build.yaml",
						DefaultText: "file",
					},
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},
					},
					&cli.BoolFlag{
						Name:  "remote",
						Usage: "Forces use remote changes",
					},
				},
			},
		},
	}

	return app.Run(args)
}

func initReporter(ctx *cli.Context) reporter.Reporter {
	impl := logrus.New()
	impl.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  time.RFC3339Nano,
		DisableTimestamp: true,
	})

	return reporter.New(
		ctx.Bool("quiet"),
		impl,
	)
}
