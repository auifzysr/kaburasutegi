package handler

import (
	"fmt"
	"os"
)

func LineChannelToken(opts ...SecretManagerOption) (string, error) {
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	if channelToken != "" {
		return channelToken, nil
	}

	// TODO: get from cmd arguments
	s := &SecretManagerConfig{
		versionID: defaultSecretVersion,
	}
	for _, opt := range opts {
		opt(s)
	}

	channelToken, err := accessSecretVersion(s)
	if err != nil {
		return "", err
	}
	return channelToken, nil
}

func LineChannelSecret() (string, error) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret == "" {
		return "", fmt.Errorf("LINE_CHANNEL_SECRET must be set")
	}
	// TOOD: get from secretmanager

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
