package quidax

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type WithdrawalsClient interface {
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
