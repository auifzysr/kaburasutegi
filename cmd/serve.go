package cmd

import (
	"log/slog"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/handler"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/service"
	"github.com/urfave/cli/v2"
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
			return serve(setup(projectID, channelSecretSecretID, channelTokenSecretID))
		},
	}
}

func setup(projectID, channelSecretSecretID, channelTokenSecretID string) (string, *service.Service) {
	slog.SetLogLoggerLevel(domain.LogLevel())

	var err error
	channelSecret, err := handler.LineChannelSecret(
		handler.WithProjectID(projectID),
		handler.WithSecretID(channelSecretSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	channelToken, err := handler.LineChannelToken(
		handler.WithProjectID(projectID),
		handler.WithSecretID(channelTokenSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	c := domain.NewCredential(channelToken, channelSecret)

	port := handler.Port()

	s := service.New(c, &domain.Journal{}, &infra.LocalRecord{})

	return port, s
}
