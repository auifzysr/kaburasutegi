package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/handler"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/service"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"github.com/urfave/cli/v2"
)

var (
	channelSecret string
	channelToken  string
)

func serve() error {
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

	cli, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		os.Exit(1)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", callbackWithAPI(cli))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	return nil
}

func serveCommand() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "serve",
		Action: func(cCtx *cli.Context) error {
			return serve()
		},
	}
}

func callbackWithAPI(cli *messaging_api.MessagingApiAPI) func(w http.ResponseWriter, req *http.Request) {
	ss := []*service.Service{
		{
			MessageHandler: domain.Register{},
			Recorder:       &infra.LocalRecord{},
		},
		{
			MessageHandler: domain.Nop{},
			Recorder:       &infra.LocalRecord{},
		},
	}

	return func(w http.ResponseWriter, req *http.Request) {
		slog.Info("/callback called...")

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

		for _, event := range cb.Events {
			slog.Debug(fmt.Sprintf("/callback called%+v...\n", event))

			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				case webhook.TextMessageContent:
					var body string
					slog.Debug(fmt.Sprintf("message: %+v", message))
					for _, s := range ss {
						if s.MessageHandler.Accept(message.Text) {
							body = s.MessageHandler.BuildMessage(message.Text)
							slog.Debug(fmt.Sprintf("body: %s", body))
							s.Recorder.Record(body)
							break
						}
					}
					if body == "" {
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
