package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

var (
	projectID             string
	channelTokenSecretID  string
	channelSecretSecretID string
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
				Name:        "channelTokenSecretID",
				Aliases:     []string{"t"},
				Value:       "",
				Usage:       "secretID of channel token",
				Destination: &channelTokenSecretID,
			},
			&cli.StringFlag{
				Name:        "channelSecretSecretID",
				Aliases:     []string{"s"},
				Value:       "",
				Usage:       "secretID of channel secret",
				Destination: &channelSecretSecretID,
			},
		},
		Commands: []*cli.Command{
			serveCommand(),
		},
	}
	return app.Run(os.Args)
}
