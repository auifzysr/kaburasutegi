package domain

import (
	"log/slog"
	"os"
)

type Credential struct {
	channelToken  string
	channelSecret string
}

func NewCredential(channelToken string, channelSecret string) *Credential {
	if channelToken == "" {
		slog.Error("channelToken must not be empty")
		os.Exit(1)
	}
	if channelSecret == "" {
		slog.Error("channelSecret must not be empty")
		os.Exit(1)
	}

	return &Credential{
		channelToken:  channelToken,
		channelSecret: channelSecret,
	}
}

func (c *Credential) GetChannelToken() string {
	return c.channelToken
}

func (c *Credential) GetChannelSecret() string {
	return c.channelSecret
}
