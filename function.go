package function

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/handler"
	"github.com/auifzysr/kaburasutegi/service"
)

func init() {
	functions.HTTP("callback", entrypoint())
}

func entrypoint() func(w http.ResponseWriter, r *http.Request) {
	_, s := setup()
	return s.Reply()
}

// cloud functions does not allow main package to reside here

func setup() (string, *service.Service) {
	slog.SetLogLoggerLevel(domain.LogLevel())
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	var err error
	channelSecret, err := handler.LineChannelSecret()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	channelToken, err := handler.LineChannelToken()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	c := domain.NewCredential(channelToken, channelSecret)

	port := handler.Port()

	s := service.New(c, service.MessageHandlersList...)

	return port, s
}
