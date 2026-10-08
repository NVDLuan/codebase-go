package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/app"
	"github.com/my-org/base-go/internal/order/domain"
)

type fakePlace struct {
	id  domain.OrderID
	err error
}

func (f fakePlace) Execute(context.Context, app.PlaceOrderCmd) (domain.OrderID, error) {
	return f.id, f.err
}

type fakeGet struct {
	view app.OrderView
	err  error
}

func (f fakeGet) Execute(context.Context, domain.OrderID) (app.OrderView, error) {
	return f.view, f.err
}

type fakeCancel struct{ err error }

func (f fakeCancel) Execute(context.Context, domain.OrderID) error { return f.err }

func serve(h *Handler, method, path, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequestWithContext(context.Background(), method, path, strings.NewReader(body)))
	return rec
}

func TestPlace(t *testing.T) {
	validBody := fmt.Sprintf(`{"customer_id":"%s","items":[{"product_id":"%s","unit_price":100,"quantity":1}]}`, uuid.New(), uuid.New())
	tests := []struct {
		name string
		uc   fakePlace
		body string
		want int
	}{
		{"created", fakePlace{id: domain.OrderID(uuid.New())}, validBody, http.StatusCreated},
		{"bad json", fakePlace{}, `{`, http.StatusBadRequest},
		{"validation", fakePlace{}, `{"items":[]}`, http.StatusUnprocessableEntity},
		{"domain error", fakePlace{err: domain.ErrEmptyOrder}, validBody, http.StatusBadRequest},
		{"infra error", fakePlace{err: errors.New("db down")}, validBody, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(NewHandler(tt.uc, fakeGet{}, fakeCancel{}), http.MethodPost, "/api/orders", tt.body)
			if rec.Code != tt.want {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, tt.want, rec.Body)
			}
		})
	}
}

func TestGet(t *testing.T) {
	id := uuid.New()
	rec := serve(NewHandler(fakePlace{}, fakeGet{view: app.OrderView{ID: id, Status: "pending"}}, fakeCancel{}),
		http.MethodGet, "/api/orders/"+id.String(), "")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), id.String()) {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
	}

	rec = serve(NewHandler(fakePlace{}, fakeGet{err: app.ErrOrderNotFound}, fakeCancel{}),
		http.MethodGet, "/api/orders/"+id.String(), "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d", rec.Code)
	}

	rec = serve(NewHandler(fakePlace{}, fakeGet{}, fakeCancel{}), http.MethodGet, "/api/orders/not-a-uuid", "")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestCancel(t *testing.T) {
	path := "/api/orders/" + uuid.NewString() + "/cancel"
	if rec := serve(NewHandler(fakePlace{}, fakeGet{}, fakeCancel{}), http.MethodPost, path, ""); rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
	err := fmt.Errorf("%w: shipped", domain.ErrNotCancellable)
	if rec := serve(NewHandler(fakePlace{}, fakeGet{}, fakeCancel{err: err}), http.MethodPost, path, ""); rec.Code != http.StatusConflict {
		t.Fatalf("status = %d", rec.Code)
	}
}
