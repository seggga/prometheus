package handler

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Home(slogger *zap.SugaredLogger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		slogger.Debug("homepage query")
		// output in plain text
		fmt.Fprintf(rw, "welcome to crop URL application")
	}
}
