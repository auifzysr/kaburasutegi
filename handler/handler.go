package handler

import (
	"fmt"
	"log/slog"
	"os"
)

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
