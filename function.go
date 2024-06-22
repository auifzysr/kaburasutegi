package function

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("Entrypoint", serve)
}

func serve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Entrypoint")
}

// cloud functions does not allow main package to reside here
