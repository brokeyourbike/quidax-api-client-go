package quidax

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type WalletData struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Balance           string    `json:"balance"`
	Locked            string    `json:"locked"`
	Staked            string    `json:"staked"`
	IsCrypto          bool      `json:"is_crypto"`
	BlockchainEnabled bool      `json:"blockchain_enabled"`
	DefaultNetwork    string    `json:"default_network"`
	DepositAddress    string    `json:"deposit_address"`
	DestinationTag    string    `json:"destination_tag"`
	Networks          []struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		DepositsEnabled  bool   `json:"deposits_enabled"`
		WithdrawsEnabled bool   `json:"withdraws_enabled"`
	} `json:"networks"`
}

func (d WalletData) GetBalance() float64 {
	f, _ := strconv.ParseFloat(d.Balance, 64)
	return f
}

func (d WalletData) GetLocked() float64 {
	f, _ := strconv.ParseFloat(d.Locked, 64)
	return f
}

func (d WalletData) GetStaked() float64 {
	f, _ := strconv.ParseFloat(d.Staked, 64)
	return f
}

type WalletResponse struct {
	Status  string     `json:"status"`
	Message string     `json:"message"`
	Data    WalletData `json:"data"`
}

type WalletsClient interface {
	FetchWallet(ctx context.Context, id uuid.UUID) (WalletResponse, error)
	FetchWallets(ctx context.Context, page int) (WalletsResponse, error)
}

func (c *client) FetchWallet(ctx context.Context, id uuid.UUID, currency string) (data WalletResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/users/%s/wallets/%s", id, currency), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

type WalletsResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []AccountData `json:"data"`
}

func (c *client) FetchWallets(ctx context.Context, id uuid.UUID) (data WalletsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}
