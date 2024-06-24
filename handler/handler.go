package handler

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/service"
)

func FunctionSetup() (string, *service.Service) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var err error
	channelSecret, err := LineChannelSecret()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	channelToken, err := LineChannelToken()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	port := Port()

	s := service.New(channelToken, channelSecret,
		&domain.Register{}, &infra.LocalRecord{})

	return port, s
}

func LocalSetup(projectID, channelSecretSecretID, channelTokenSecretID string) (string, *service.Service) {
	if env := Env(); env == "local" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	var err error
	channelSecret, err := LineChannelSecret(
		WithProjectID(projectID),
		WithSecretID(channelSecretSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	channelToken, err := LineChannelToken(
		WithProjectID(projectID),
		WithSecretID(channelTokenSecretID),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	port := Port()

	s := service.New(channelToken, channelSecret,
		&domain.Register{}, &infra.LocalRecord{})

	return port, s
}

func LineChannelToken(opts ...SecretManagerParam) (string, error) {
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	if channelToken != "" {
		return channelToken, nil
	}

	s := &SecretManagerConfig{
		versionID: defaultSecretVersion,
	}
	for _, opt := range opts {
		opt(s)
	}

	slog.Debug("rerieving channel token  from secret manager...")

	channelToken, err := accessSecretVersion(s)
	if err != nil {
		return "", fmt.Errorf("failed to get channel token: %w", err)
	}
	return channelToken, nil
}

func LineChannelSecret(opts ...SecretManagerParam) (string, error) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret != "" {
		return channelSecret, nil
	}

	s := &SecretManagerConfig{
		versionID: defaultSecretVersion,
	}
	for _, opt := range opts {
		opt(s)
	}

	slog.Debug("rerieving channel secret from secret manager...")

	channelSecret, err := accessSecretVersion(s)
	if err != nil {
		return "", fmt.Errorf("failed to get channel secret: %w", err)
	}

	return channelSecret, nil
}

const (
	DEFAULT_ENV  = "local"
	DEFAULT_PORT = "3000"
)

func Env() string {
	env := os.Getenv("ENV")
	if env == "" {
		return DEFAULT_ENV
	}
	return env
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		return DEFAULT_PORT
	}
	return port
}
