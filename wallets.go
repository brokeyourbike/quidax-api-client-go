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
	FetchWallet(ctx context.Context, id uuid.UUID, currency string) (WalletResponse, error)
	FetchWallets(ctx context.Context, id uuid.UUID) (WalletsResponse, error)
	FetchWalletAddress(ctx context.Context, id uuid.UUID, currency string) (data WalletAddressResponse, err error)
	FetchWalletAddresses(ctx context.Context, id uuid.UUID, currency string) (data WalletAddressesResponse, err error)
	RequestWalletAddress(ctx context.Context, id uuid.UUID, currency, network string) (data WalletAddressResponse, err error)
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

type WalletAddressData struct {
	ID             uuid.UUID `json:"id"`
	Currency       string    `json:"currency"`
	Address        string    `json:"address"`
	Network        string    `json:"network"`
	Reference      string    `json:"reference"`
	DestinationTag string    `json:"destination_tag"`
}

type WalletAddressResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    WalletAddressData `json:"data"`
}

func (c *client) FetchWalletAddress(ctx context.Context, id uuid.UUID, currency string) (data WalletAddressResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/users/%s/wallets/%s/address", id, currency), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

type WalletAddressesResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []WalletAddressData `json:"data"`
}

func (c *client) FetchWalletAddresses(ctx context.Context, id uuid.UUID, currency string) (data WalletAddressesResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/users/%s/wallets/%s/addresses", id, currency), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

func (c *client) RequestWalletAddress(ctx context.Context, id uuid.UUID, currency, network string) (data WalletAddressResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/v1/users/%s/wallets/%s/address", id, currency), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	if network != "" {
		req.AddQueryParam("network", network)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusCreated)
	return data, c.do(ctx, req)
}
