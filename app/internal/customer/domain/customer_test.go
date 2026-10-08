package domain

import (
	"errors"
	"testing"
)

func TestNewCustomer(t *testing.T) {
	if _, err := NewCustomer("A@Example.com", "An"); err != nil {
		t.Fatal(err)
	}
	if _, err := NewCustomer("bad", "An"); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("err = %v", err)
	}
	if _, err := NewCustomer("a@example.com", " "); !errors.Is(err, ErrEmptyName) {
		t.Fatalf("err = %v", err)
	}
}
