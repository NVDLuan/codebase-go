// Package postgres hiện thực app.Repository bằng pgx + code sqlc sinh ra (thư mục db/).
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/my-org/base-go/internal/order/adapter/postgres/db"
	"github.com/my-org/base-go/internal/order/domain"
)

type Repository struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool, q: db.New(pool)}
}

func (r *Repository) Save(ctx context.Context, o *domain.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.q.WithTx(tx)
	id := uuid.UUID(o.ID())
	if err := q.UpsertOrder(ctx, db.UpsertOrderParams{
		ID:         id,
		CustomerID: o.CustomerID(),
		Status:     string(o.Status()),
	}); err != nil {
		return fmt.Errorf("upsert order: %w", err)
	}
	if err := q.DeleteOrderItems(ctx, id); err != nil {
		return fmt.Errorf("delete items: %w", err)
	}
	for i, it := range o.Items() {
		if err := q.InsertOrderItem(ctx, db.InsertOrderItemParams{
			OrderID:   id,
			Position:  int32(i),
			ProductID: it.ProductID,
			UnitPrice: it.UnitPrice.Amount,
			Currency:  it.UnitPrice.Currency,
			Quantity:  int32(it.Quantity),
		}); err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *Repository) FindByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	row, err := r.q.GetOrder(ctx, uuid.UUID(id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}
	rows, err := r.q.ListOrderItems(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	items := make([]domain.OrderItem, 0, len(rows))
	for _, it := range rows {
		items = append(items, domain.OrderItem{
			ProductID: it.ProductID,
			UnitPrice: domain.Money{Amount: it.UnitPrice, Currency: it.Currency},
			Quantity:  int(it.Quantity),
		})
	}
	return domain.Rehydrate(domain.OrderID(row.ID), row.CustomerID, items, domain.Status(row.Status)), nil
}
