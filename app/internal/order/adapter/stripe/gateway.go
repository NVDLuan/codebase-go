// Package stripe hiện thực app.PaymentGateway bằng Stripe PaymentIntents API.
package stripe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/my-org/base-go/internal/order/domain"
)

const defaultBaseURL = "https://api.stripe.com"

type Gateway struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewGateway(apiKey string) *Gateway {
	return &Gateway{apiKey: apiKey, baseURL: defaultBaseURL, client: &http.Client{Timeout: 10 * time.Second}}
}

func (g *Gateway) Charge(ctx context.Context, customerID uuid.UUID, amount domain.Money) error {
	form := url.Values{
		"amount":                {strconv.FormatInt(amount.Amount, 10)},
		"currency":              {strings.ToLower(amount.Currency)},
		"metadata[customer_id]": {customerID.String()},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.baseURL+"/v1/payment_intents", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(g.apiKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("stripe request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()
	if res.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		return fmt.Errorf("stripe: status %d: %s", res.StatusCode, body)
	}
	return nil
}
