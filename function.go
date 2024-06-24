package function

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/auifzysr/kaburasutegi/handler"
)

func init() {
	_, s := handler.FunctionSetup()
	functions.HTTP("Entrypoint", s.Respond())
	functions.HTTP("healthcheck", healthcheck)
}

// cloud functions does not allow main package to reside here
func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}
