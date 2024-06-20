package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"regexp"

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

type messageHandler interface {
	accept(text string) bool
	buildMessage(text string) string
}

type nop struct{}

func (s nop) accept(text string) bool {
	return true
}

func (s nop) buildMessage(text string) string {
	return text
}

type register struct{}

func (s register) accept(text string) bool {
	match, err := regexp.MatchString(`^(?:[01][0-9]|2[0-3])[0-5][0-9] .*$`, text)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to match: %s", text))
		return false
	}
	return match
}

func (s register) buildMessage(text string) string {
	return fmt.Sprintf("registered succcessfully: %s", text)
}

type recorder interface {
	record(text string) error
	recordAt(id int, text string) error
	readAt(id int) (string, error)
}

type localRecord struct {
	records map[int]string
}

func (s *localRecord) record(text string) error {
	return s.recordAt(len(s.records), text)
}

func (s *localRecord) recordAt(id int, text string) error {
	if s.records == nil {
		s.records = make(map[int]string)
	}
	s.records[id] = text
	return nil
}

func (s *localRecord) readAt(id int) (string, error) {
	if s.records == nil {
		s.records = make(map[int]string)
	}
	if text, ok := s.records[id]; ok {
		return text, nil
	}
	return "", fmt.Errorf("not found")
}

func callbackWithAPI(cli *messaging_api.MessagingApiAPI) func(w http.ResponseWriter, req *http.Request) {
	messageHandlers := []messageHandler{
		register{},
		nop{},
	}

	recorder := &localRecord{}

	return func(w http.ResponseWriter, req *http.Request) {
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
					for _, builder := range messageHandlers {
						if builder.accept(message.Text) {
							slog.Debug(fmt.Sprintf("response builder: %s", reflect.TypeOf(builder).Name()))
							body = builder.buildMessage(message.Text)
							recorder.record(body)
							slog.Debug(fmt.Sprintf("records: %+v", recorder.records))
							break
						}
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
