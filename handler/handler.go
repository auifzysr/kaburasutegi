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
	slog.SetLogLoggerLevel(LogLevel())
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

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

	c := domain.NewCredential(channelToken, channelSecret)

	port := Port()

	s := service.New(c, &domain.Journal{}, &infra.LocalRecord{})

	return port, s
}

func LogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func LocalSetup(projectID, channelSecretSecretID, channelTokenSecretID string) (string, *service.Service) {
	slog.SetLogLoggerLevel(LogLevel())

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

	c := domain.NewCredential(channelToken, channelSecret)

	port := Port()

	s := service.New(c, &domain.Journal{}, &infra.LocalRecord{})

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
