package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

var (
	projectID string
	secretID  string
	versionID string
)

func Run() error {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "projectID",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "projectID",
				Destination: &projectID,
			},
			&cli.StringFlag{
				Name:        "secretID",
				Aliases:     []string{"s"},
				Value:       "",
				Usage:       "secretID",
				Destination: &secretID,
			},
			&cli.StringFlag{
				Name:        "versionID",
				Aliases:     []string{"v"},
				Value:       "",
				Usage:       "versionID",
				Destination: &versionID,
			},
		},
		Commands: []*cli.Command{
			serveCommand(),
		},
	}
	return app.Run(os.Args)
}
