package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/repository"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Service struct {
	channelToken  string
	channelSecret string

	messageHandler domain.MessageBuilder
	recorder       repository.Recorder
}

func New(channelToken string, channelSecret string, opts ...interface{}) *Service {
	if channelToken == "" {
		slog.Error("channelToken is empty")
		os.Exit(1)
	}
	if channelSecret == "" {
		slog.Error("channelSecret is empty")
		os.Exit(1)
	}

	s := &Service{}
	s.channelSecret = channelSecret
	s.channelToken = channelToken

	for _, opt := range opts {
		switch opt.(type) {
		case domain.MessageBuilder:
			s.messageHandler = opt.(domain.MessageBuilder)
		case repository.Recorder:
			s.recorder = opt.(repository.Recorder)
		default:
			slog.Error(fmt.Sprintf("Unsupported option: %T\n", opt))
			os.Exit(1)
		}
	}
	if s.messageHandler == nil {
		s.messageHandler = domain.Nop{}
	}
	if s.recorder == nil {
		s.recorder = &infra.LocalRecord{}
	}

	return s
}

func (s *Service) Respond() func(w http.ResponseWriter, req *http.Request) {
	cli, err := messaging_api.NewMessagingApiAPI(s.channelToken)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		slog.Info("/callback called...")

		cb, err := webhook.ParseRequest(s.channelSecret, req)
		if err != nil {
			slog.Error(fmt.Sprintf("Cannot parse request: %+v\n", err))
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range cb.Events {
			slog.Debug(fmt.Sprintf("/callback called%+v...\n", event))

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					var body string
					slog.Debug(fmt.Sprintf("message: %+v", message))
					if s.messageHandler.Accept(message.Text) {
						body = s.messageHandler.BuildMessage(message.Text)
						slog.Debug(fmt.Sprintf("body: %s", body))
						s.recorder.Record(body)
					} else {
						body = "error: no such handler"
					}

					if _, err = cli.ReplyMessage(
						&messaging_api.ReplyMessageRequest{
							ReplyToken: e.ReplyToken,
							Messages: []messaging_api.MessageInterface{
								messaging_api.TextMessage{
									Text: body,
								},
							},
						},
					); err != nil {
						slog.Error(fmt.Sprintf("%s", err))
						w.WriteHeader(500)
					} else {
						slog.Debug("Sent text reply.")
					}
				default:
					slog.Error(fmt.Sprintf("Unsupported message content: %T\n", e.Message))
					w.WriteHeader(400)
				}
			default:
				slog.Error(fmt.Sprintf("Unsupported message: %T\n", event))
				w.WriteHeader(400)
			}
		}
	}
}
