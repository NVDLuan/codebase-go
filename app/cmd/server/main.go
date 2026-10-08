// Composition root: đọc config, tạo adapter, ráp use case, chạy HTTP.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	orderkafka "github.com/my-org/base-go/internal/order/adapter/kafka"
	orderlocal "github.com/my-org/base-go/internal/order/adapter/local"
	orderpg "github.com/my-org/base-go/internal/order/adapter/postgres"
	orderstripe "github.com/my-org/base-go/internal/order/adapter/stripe"
	orderapp "github.com/my-org/base-go/internal/order/app"
	orderhttp "github.com/my-org/base-go/internal/order/transport/http"
	"github.com/my-org/base-go/internal/platform/config"
	"github.com/my-org/base-go/internal/platform/database"
	"github.com/my-org/base-go/internal/platform/httpserver"
	"github.com/my-org/base-go/internal/platform/observability"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := observability.NewLogger(cfg.LogLevel).With("env", cfg.AppEnv)
	slog.SetDefault(log)

	// ---------- Adapters ----------
	var (
		repo      orderapp.Repository
		payment   orderapp.PaymentGateway
		publisher orderapp.Publisher
		pool      *pgxpool.Pool
	)

	if cfg.DatabaseURL != "" {
		pool, err = database.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		if err := database.Migrate(ctx, pool); err != nil {
			return err
		}
		repo = orderpg.NewRepository(pool)
	} else {
		log.Warn("DATABASE_URL not set: using in-memory repository")
		repo = orderlocal.NewRepository()
	}

	if cfg.StripeKey != "" {
		payment = orderstripe.NewGateway(cfg.StripeKey)
	} else {
		log.Warn("STRIPE_API_KEY not set: using local payment gateway")
		payment = orderlocal.PaymentGateway{Log: log}
	}

	if len(cfg.KafkaBrokers) > 0 {
		p := orderkafka.NewPublisher(cfg.KafkaBrokers, cfg.KafkaTopic)
		defer func() { _ = p.Close() }()
		publisher = p
	} else {
		log.Warn("KAFKA_BROKERS not set: using log publisher")
		publisher = orderlocal.Publisher{Log: log}
	}

	// ---------- Use cases ----------
	placeOrder := orderapp.NewPlaceOrder(repo, payment, publisher)
	getOrder := orderapp.NewGetOrder(repo)
	cancelOrder := orderapp.NewCancelOrder(repo)

	// ---------- Transport ----------
	mux := http.NewServeMux()
	if pool != nil {
		httpserver.RegisterHealth(mux, pool)
	} else {
		httpserver.RegisterHealth(mux)
	}
	orderhttp.NewHandler(placeOrder, getOrder, cancelOrder).Register(mux)

	handler := httpserver.Chain(mux, httpserver.Recover(log), httpserver.RequestID, httpserver.Logging(log))
	return httpserver.Run(ctx, cfg.HTTPAddr, handler, cfg.ShutdownTimeout)
}
