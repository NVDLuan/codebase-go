// Package domain của feature customer. Feature mới: làm theo cấu trúc của internal/order.
package domain

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrEmptyName    = errors.New("name is required")
)

type CustomerID uuid.UUID

type Customer struct {
	id    CustomerID
	email string
	name  string
}

func NewCustomer(email, name string) (*Customer, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrInvalidEmail
	}
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}
	return &Customer{id: CustomerID(uuid.New()), email: strings.ToLower(email), name: name}, nil
}

func (c *Customer) ID() CustomerID { return c.id }
func (c *Customer) Email() string  { return c.email }
func (c *Customer) Name() string   { return c.name }
