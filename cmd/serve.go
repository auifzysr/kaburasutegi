package cmd

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/auifzysr/kaburasutegi/handler"
	"github.com/auifzysr/kaburasutegi/service"
	"github.com/urfave/cli/v2"
)

var (
	channelSecret string
	channelToken  string
)

const (
	defaultHostname = "0.0.0.0"
)

func serve(port string, s *service.Service) error {
	functions.HTTP("callback", s.Respond())

	return funcframework.StartHostPort(defaultHostname, port)
}

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "serve",
		Action: func(cCtx *cli.Context) error {
			return serve(handler.LocalSetup(projectID, channelSecretSecretID, channelTokenSecretID))
		},
	}
}
