package v1alpha1

import (
	"fmt"
	"log/slog"
	"net/http"

	oapimw "github.com/oapi-codegen/nethttp-middleware"
	serverapiv1alpha1 "github.com/spacechunks/control-plane/api/server/v1alpha1"
	apierr "github.com/spacechunks/control-plane/internal/error"
	"github.com/spacechunks/control-plane/internal/httphelper"
)

var _ serverapiv1alpha1.StrictServerInterface = (*Handler)(nil)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func AddToMux(logger *slog.Logger, h *Handler, mux *http.ServeMux) error {
	spec, err := serverapiv1alpha1.GetSpec()
	if err != nil {
		return fmt.Errorf("spec: %w", err)
	}

	validator := oapimw.OapiRequestValidatorWithOptions(spec, &oapimw.Options{
		SilenceServersWarning: true,
		DoNotValidateServers:  true,
	})

	serverapiv1alpha1.HandlerFromMux(
		serverapiv1alpha1.NewStrictHandlerWithOptions(h, nil, serverapiv1alpha1.StrictHTTPServerOptions{
			ResponseErrorHandlerFunc: apierr.OAPIRequestErrorHandler(logger),
		}),
		httphelper.WithMiddleware(mux, validator),
	)

	return nil
}
