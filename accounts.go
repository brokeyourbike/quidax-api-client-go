package quidax

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const ParentAccountID = "me"

type AccountData struct {
	ID          uuid.UUID `json:"id"`
	SN          string    `json:"sn"`
	Email       string    `json:"email"`
	Reference   string    `json:"reference"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DisplayName string    `json:"display_name"`
}

type AccountResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    AccountData `json:"data"`
}

type AccountsClient interface {
	FetchParentAccount(ctx context.Context) (AccountResponse, error)
	FetchAccount(ctx context.Context, id uuid.UUID) (AccountResponse, error)
	FetchAccounts(ctx context.Context, page int) (AccountsResponse, error)
	CreateAccount(ctx context.Context, payload CreateAccountPayload) (AccountResponse, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, payload UpdateAccountPayload) (AccountResponse, error)
}

func (c *client) FetchParentAccount(ctx context.Context) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/users/%s", ParentAccountID), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

func (c *client) FetchAccount(ctx context.Context, id uuid.UUID) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/v1/users/%s", id), nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

type AccountsResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []AccountData `json:"data"`
}

func (c *client) FetchAccounts(ctx context.Context, page int) (data AccountsResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users", nil)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddQueryParam("per_page", strconv.Itoa(c.perPage))
	req.AddQueryParam("page", strconv.Itoa(page))
	req.DecodeTo(&data)
	req.ExpectStatus(http.StatusOK)
	return data, c.do(ctx, req)
}

type CreateAccountPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *client) CreateAccount(ctx context.Context, payload CreateAccountPayload) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/users", payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusCreated)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}

type UpdateAccountPayload struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

func (c *client) UpdateAccount(ctx context.Context, id uuid.UUID, payload UpdateAccountPayload) (data AccountResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/v1/users/%s", id), payload)
	if err != nil {
		return data, fmt.Errorf("failed to create request: %w", err)
	}

	req.ExpectStatus(http.StatusOK)
	req.DecodeTo(&data)
	return data, c.do(ctx, req)
}
