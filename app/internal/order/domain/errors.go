package domain

import "errors"

var (
	ErrEmptyOrder      = errors.New("order must have at least one item")
	ErrMixedCurrency   = errors.New("all items must use the same currency")
	ErrInvalidQuantity = errors.New("quantity must be greater than 0")
	ErrNegativePrice   = errors.New("unit price must not be negative")
	ErrNotCancellable  = errors.New("order cannot be cancelled")
)
