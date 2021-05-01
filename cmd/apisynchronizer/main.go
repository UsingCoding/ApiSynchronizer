package main

import (
	"github.com/urfave/cli/v2"
	stdlog "log"
	"os"
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
		Name:    appID,
		Version: version,
		Commands: []*cli.Command{
			{
				Name:      "resolve",
				UsageText: "Resolves api-files for service",
				Action:    executeResolve,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Required: true,
					},
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
					},
				},
			},
		},
	}

	return app.Run(args)
}
