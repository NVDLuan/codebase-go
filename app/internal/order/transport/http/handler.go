// Package http là lớp transport (presentation) của feature order.
package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/app"
	"github.com/my-org/base-go/internal/order/domain"
	"github.com/my-org/base-go/internal/platform/httpserver"
)

// Handler chỉ cần đúng method nó dùng → interface nhỏ, dễ test.
type placeOrderUC interface {
	Execute(ctx context.Context, cmd app.PlaceOrderCmd) (domain.OrderID, error)
}

type getOrderUC interface {
	Execute(ctx context.Context, id domain.OrderID) (app.OrderView, error)
}

type cancelOrderUC interface {
	Execute(ctx context.Context, id domain.OrderID) error
}

type Handler struct {
	placeOrder  placeOrderUC
	getOrder    getOrderUC
	cancelOrder cancelOrderUC
}

func NewHandler(po placeOrderUC, gt getOrderUC, co cancelOrderUC) *Handler {
	return &Handler{placeOrder: po, getOrder: gt, cancelOrder: co}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/orders", h.place)
	mux.HandleFunc("GET /api/orders/{id}", h.get)
	mux.HandleFunc("POST /api/orders/{id}/cancel", h.cancel)
}

const maxBodyBytes = 1 << 20

func (h *Handler) place(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&req); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := req.Validate(); err != nil {
		httpserver.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	id, err := h.placeOrder.Execute(r.Context(), req.toCmd())
	if err != nil {
		writeError(w, err)
		return
	}
	httpserver.WriteJSON(w, http.StatusCreated, orderCreatedResponse{ID: uuid.UUID(id)})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	view, err := h.getOrder.Execute(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	httpserver.WriteJSON(w, http.StatusOK, fromView(view))
}

func (h *Handler) cancel(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := h.cancelOrder.Execute(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func pathID(w http.ResponseWriter, r *http.Request) (domain.OrderID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, "invalid order id")
		return domain.OrderID{}, false
	}
	return domain.OrderID(id), true
}
