package quidax

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type WithdrawalsClient interface {
	FetchWithdrawalFees(ctx context.Context, currency, network string) (FeesResponse, error)
	CreateWithdrawal(ctx context.Context, userID uuid.UUID, payload CreateWithdrawalPayload) (WithdrawalResponse, error)
}

type WithdrawalData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type WithdrawalResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    WithdrawalData `json:"data"`
}

type CreateWithdrawalPayload struct {
	Currency        string `json:"currency"`
	Amount          string `json:"amount"`
	TransactionNote string `json:"transaction_note"`
	Narration       string `json:"narration"`
	FundUID         string `json:"fund_uid"`
	Network         string `json:"network"`
	Reference       string `json:"reference"`
}

func (c *client) CreateWithdrawal(ctx context.Context, userID uuid.UUID, payload CreateWithdrawalPayload) (data WithdrawalResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/users/%s/withdraws", userID), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type FeeData struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type OneFeeData struct {
	Type  string  `json:"type"`
	Value float64 `json:"fee"`
}

type FeesResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (r FeesResponse) GetFees() []FeeData {
	fees := make([]FeeData, 0)

	var one OneFeeData
	if err := json.Unmarshal(r.Data, &one); err == nil {
		fees = append(fees, FeeData{Type: one.Type, Value: one.Value})
	} else {
		var multi []FeeData
		if err := json.Unmarshal(r.Data, &multi); err == nil {
			fees = append(fees, multi...)
		}
	}

	return fees
}

func (c *client) FetchWithdrawalFees(ctx context.Context, currency, network string) (data FeesResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/fee", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("currency", strings.ToLower(currency))

	if network != "" {
		req.AddQueryParam("network", strings.ToLower(network))
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}
