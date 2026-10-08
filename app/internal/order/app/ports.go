package app

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

// ErrOrderNotFound do use case trả về khi Repository không tìm thấy đơn.
var ErrOrderNotFound = errors.New("order not found")

// Interface đặt ở phía dùng; adapter thỏa mãn ngầm định.
type Repository interface {
	// Save tạo mới hoặc cập nhật đơn hàng (kèm item) trong 1 transaction.
	Save(ctx context.Context, o *domain.Order) error
	// FindByID trả (nil, nil) khi không tìm thấy.
	FindByID(ctx context.Context, id domain.OrderID) (*domain.Order, error)
}

type PaymentGateway interface {
	Charge(ctx context.Context, customerID uuid.UUID, amount domain.Money) error
}

type Publisher interface {
	Publish(ctx context.Context, events []any) error
}
