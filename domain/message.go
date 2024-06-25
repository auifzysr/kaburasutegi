package domain

import (
	"fmt"
	"log/slog"
	"regexp"
)

type MessageBuilder interface {
	Accept(text string) bool
	BuildReply(text string) string
}

type Nop struct{}

func (s Nop) Accept(text string) bool {
	return true
}

func (s Nop) BuildReply(text string) string {
	return text
}

type Journal struct{}

func (s Journal) Accept(text string) bool {
	match, err := regexp.MatchString(`^(?:[01][0-9]|2[0-3])[0-5][0-9] .*$`, text)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to match: %s", text))
		return false
	}
	return match
}

func (s Journal) BuildReply(text string) string {
	return fmt.Sprintf("registered succcessfully: %s", text)
}
