package function

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/auifzysr/kaburasutegi/handler"
)

func init() {
	functions.HTTP("callback", entrypoint())
}

func entrypoint() func(w http.ResponseWriter, r *http.Request) {
	_, s := handler.FunctionSetup()
	return s.Respond()
}

// cloud functions does not allow main package to reside here
