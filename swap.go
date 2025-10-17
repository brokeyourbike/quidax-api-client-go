package quidax

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SwapClient interface {
	Quote(ctx context.Context, userID uuid.UUID, payload QuotePayload) (QuoteResponse, error)
	ConfirmQuote(ctx context.Context, userID, quoteID uuid.UUID) error
}

type QuotePayload struct {
	FromCurrency string `json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	FromAmount   string `json:"from_amount,omitempty"`
	ToAmount     string `json:"to_amount,omitempty"`
}

type QuoteData struct {
	ID           uuid.UUID `json:"id"`
	FromCurrency string    `json:"from_currency"`
	ToCurrency   string    `json:"to_currency"`
	FromAmount   string    `json:"from_amount"`
	ToAmount     string    `json:"to_amount"`
	Confirmed    bool      `json:"confirmed"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (d QuoteData) GetFromAmount() float64 {
	f, _ := strconv.ParseFloat(d.FromAmount, 64)
	return f
}

func (d QuoteData) GetToAmount() float64 {
	f, _ := strconv.ParseFloat(d.ToAmount, 64)
	return f
}

type QuoteResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    QuoteData `json:"data"`
}

func (c *client) Quote(ctx context.Context, userID uuid.UUID, payload QuotePayload) (data QuoteResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/users/%s/swap_quotation", userID), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

func (c *client) ConfirmQuote(ctx context.Context, userID, quoteID uuid.UUID) error {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/users/%s/swap_quotation/%s/confirm", userID, quoteID), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	return c.do(ctx, req)
}
