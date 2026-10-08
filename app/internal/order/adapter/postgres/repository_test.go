//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/my-org/base-go/internal/order/domain"
	"github.com/my-org/base-go/internal/platform/database"
)

// Chạy: go test -tags=integration ./internal/order/adapter/postgres/... (cần Docker)
func TestRepository(t *testing.T) {
	ctx := context.Background()
	ctr, err := tcpostgres.Run(ctx, "postgres:16",
		tcpostgres.WithDatabase("shop"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(60*time.Second)))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })

	url, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	pool, err := database.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}

	repo := NewRepository(pool)
	o, err := domain.NewOrder(uuid.New(), []domain.OrderItem{
		{ProductID: uuid.New(), UnitPrice: domain.Money{Amount: 100, Currency: "VND"}, Quantity: 2},
		{ProductID: uuid.New(), UnitPrice: domain.Money{Amount: 30, Currency: "VND"}, Quantity: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, o); err != nil {
		t.Fatal(err)
	}

	got, err := repo.FindByID(ctx, o.ID())
	if err != nil || got == nil {
		t.Fatalf("find: %v %v", got, err)
	}
	if got.Total() != o.Total() || got.Status() != domain.StatusPending || len(got.Items()) != 2 {
		t.Fatalf("got %+v", got)
	}

	// Cập nhật trạng thái
	if err := got.Cancel(); err != nil {
		t.Fatal(err)
	}
	if err := repo.Save(ctx, got); err != nil {
		t.Fatal(err)
	}
	again, _ := repo.FindByID(ctx, o.ID())
	if again.Status() != domain.StatusCancelled || len(again.Items()) != 2 {
		t.Fatalf("after cancel: %+v", again)
	}

	missing, err := repo.FindByID(ctx, domain.OrderID(uuid.New()))
	if err != nil || missing != nil {
		t.Fatalf("missing: %v %v", missing, err)
	}
}
