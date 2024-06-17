package main

import (
	"errors"
	"fmt"
	"log"
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
	channelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	if channelSecret == "" {
		log.Fatal("LINE_CHANNEL_SECRET must be set")
	}
	channelToken = os.Getenv("LINE_CHANNEL_TOKEN")
	if channelToken == "" {
		log.Fatal("LINE_CHANNEL_TOKEN must be set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	cli, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", callbackWithAPI(cli))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type callbackFunc func(w http.ResponseWriter, req *http.Request)

var messages []string

func callbackWithAPI(cli *messaging_api.MessagingApiAPI) callbackFunc {
	return callbackFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("/callback called...")

		cb, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			log.Printf("Cannot parse request: %+v\n", err)
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		log.Println("Handling events...")
		for _, event := range cb.Events {
			log.Printf("/callback called%+v...\n", event)

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
						log.Print(err)
					} else {
						log.Println("Sent text reply.")
					}
				default:
					log.Printf("Unsupported message content: %T\n", e.Message)
				}
			default:
				log.Printf("Unsupported message: %T\n", event)
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
