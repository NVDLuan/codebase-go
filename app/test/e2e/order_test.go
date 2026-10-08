//go:build e2e

// E2E: ráp app thật với Postgres (testcontainers) và gọi qua HTTP.
// Chạy: go test -tags=e2e ./test/e2e/... (cần Docker)
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	orderlocal "github.com/my-org/base-go/internal/order/adapter/local"
	orderpg "github.com/my-org/base-go/internal/order/adapter/postgres"
	orderapp "github.com/my-org/base-go/internal/order/app"
	orderhttp "github.com/my-org/base-go/internal/order/transport/http"
	"github.com/my-org/base-go/internal/platform/database"
	"github.com/my-org/base-go/internal/platform/httpserver"
)

func TestOrderFlow(t *testing.T) {
	ctx := context.Background()
	ctr, err := tcpostgres.Run(ctx, "postgres:16",
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(60*time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })
	url, _ := ctr.ConnectionString(ctx, "sslmode=disable")
	pool, err := database.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	repo := orderpg.NewRepository(pool)
	mux := http.NewServeMux()
	httpserver.RegisterHealth(mux, pool)
	orderhttp.NewHandler(
		orderapp.NewPlaceOrder(repo, orderlocal.PaymentGateway{Log: log}, orderlocal.Publisher{Log: log}),
		orderapp.NewGetOrder(repo),
		orderapp.NewCancelOrder(repo),
	).Register(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	// Đặt hàng
	body := fmt.Sprintf(`{"customer_id":"%s","items":[{"product_id":"%s","unit_price":1000,"quantity":3}]}`, uuid.New(), uuid.New())
	var created struct{ ID string }
	if code := call(t, http.MethodPost, srv.URL+"/api/orders", body, &created); code != http.StatusCreated {
		t.Fatalf("place: %d", code)
	}

	// Xem
	var got struct {
		Status string
		Total  int64
	}
	if code := call(t, http.MethodGet, srv.URL+"/api/orders/"+created.ID, "", &got); code != http.StatusOK || got.Total != 3000 || got.Status != "pending" {
		t.Fatalf("get: %d %+v", code, got)
	}

	// Huỷ
	if code := call(t, http.MethodPost, srv.URL+"/api/orders/"+created.ID+"/cancel", "", nil); code != http.StatusNoContent {
		t.Fatalf("cancel: %d", code)
	}

	// Readiness
	if code := call(t, http.MethodGet, srv.URL+"/health/ready", "", nil); code != http.StatusOK {
		t.Fatalf("ready: %d", code)
	}
}

// call gửi request, decode JSON vào out (nếu khác nil) và trả status code.
func call(t *testing.T, method, url, body string, out any) int {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), method, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatal(err)
		}
	}
	return res.StatusCode
}
