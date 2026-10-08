package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/my-org/base-go/internal/order/app"
	"github.com/my-org/base-go/internal/order/domain"
	"github.com/my-org/base-go/internal/platform/httpserver"
)

// writeError là nơi DUY NHẤT map lỗi của feature order sang HTTP status.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, app.ErrOrderNotFound):
		httpserver.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrNotCancellable):
		httpserver.WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrEmptyOrder),
		errors.Is(err, domain.ErrMixedCurrency),
		errors.Is(err, domain.ErrInvalidQuantity),
		errors.Is(err, domain.ErrNegativePrice):
		httpserver.WriteError(w, http.StatusBadRequest, err.Error())
	default:
		slog.Error("order: internal error", "err", err)
		httpserver.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}
