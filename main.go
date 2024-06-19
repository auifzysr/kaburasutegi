package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

var (
	channelSecret string
	channelToken  string
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	channelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret == "" {
		slog.Error("LINE_CHANNEL_SECRET must be set")
	}
	channelToken = os.Getenv("LINE_CHANNEL_TOKEN")
	if channelToken == "" {
		slog.Error("LINE_CHANNEL_TOKEN must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	cli, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", callbackWithAPI(cli))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error(fmt.Sprintf("%s", err))
	}
}

type callbackFunc func(w http.ResponseWriter, req *http.Request)

var messages []string

func callbackWithAPI(cli *messaging_api.MessagingApiAPI) callbackFunc {
	return callbackFunc(func(w http.ResponseWriter, req *http.Request) {
		slog.Debug("/callback called...")

		cb, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			slog.Error(fmt.Sprintf("Cannot parse request: %+v\n", err))
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		slog.Debug("Handling events...")
		for _, event := range cb.Events {
			slog.Debug(fmt.Sprintf("/callback called%+v...\n", event))

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					var body string
					if message.Text == "一覧" {
						body = listMessage(messages)
					} else {
						body = buildMessage(message.Text)
						messages = append(messages, message.Text)
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
	})
}

func buildMessage(text string) string {
	return fmt.Sprintf("\"%s\" 記録しました", text)
}

func listMessage(messages []string) string {
	var body string
	for _, message := range messages {
		body += message + "\n"
	}
	return fmt.Sprintf("記録一覧\n%s", body)
}
