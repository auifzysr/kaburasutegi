package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Service struct {
	credential      *domain.Credential
	messageHandlers []*MessageHandler
}

func New(credential *domain.Credential, handlers ...*MessageHandler) *Service {
	s := &Service{
		credential:      credential,
		messageHandlers: handlers,
	}

	return s
}

func (s *Service) Reply() func(w http.ResponseWriter, req *http.Request) {
	cli, err := messaging_api.NewMessagingApiAPI(s.credential.GetChannelToken())
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		slog.Info("/callback called...")

		cb, err := webhook.ParseRequest(s.credential.GetChannelSecret(), req)
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
					for _, h := range s.messageHandlers {
						if h.messageBuilder.Accept(message.Text) {
							body = h.messageBuilder.BuildReply(message.Text)
							slog.Debug(fmt.Sprintf("body: %s", body))
							h.recorder.Record(body)
							break
						}
						body = "error: no such handler"
					}

					slog.Debug(fmt.Sprintf("body: %s", body))
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
