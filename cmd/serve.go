package cmd

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/handler"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/service"
	"github.com/urfave/cli/v2"
)

var (
	channelSecret string
	channelToken  string
)

func serve(port string, s *service.Service) error {
	http.HandleFunc("/callback", s.Respond())

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return nil
}

func setup() (string, *service.Service) {
	if env := handler.Env(); env == "local" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	var err error
	channelSecret, err = handler.LineChannelSecret(
		handler.WithProjectID(projectID),
		handler.WithSecretID(channelSecretSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	channelToken, err = handler.LineChannelToken(
		handler.WithProjectID(projectID),
		handler.WithSecretID(channelTokenSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	port := handler.Port()

	s := service.New(channelToken, channelSecret,
		&domain.Register{}, &infra.LocalRecord{})

	return port, s
}

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "serve",
		Action: func(cCtx *cli.Context) error {
			return serve(setup())
		},
	}
}
