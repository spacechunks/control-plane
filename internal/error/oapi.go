package error

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func OAPIRequestErrorHandler(logger *slog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		ret := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "internal",
			Message: "an internal error occurred",
		}

		logger.Error("internal error occurred", "err", err)
		_ = json.NewEncoder(w).Encode(ret)
	}
}
