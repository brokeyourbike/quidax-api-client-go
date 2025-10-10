package quidax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

const defaultBaseURL = "https://app.quidax.io/api"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	AccountsClient
	WalletsClient
	WithdrawalsClient
	SwapClient
}

var _ Client = (*client)(nil)

type client struct {
	httpClient HttpClient
	logger     *logrus.Logger
	baseURL    string
	token      string
	perPage    int
}

// ClientOption is a function that configures a Client.
type ClientOption func(*client)

// WithHTTPClient sets the HTTP client for the ClearBank API client.
func WithHTTPClient(c HttpClient) ClientOption {
	return func(target *client) {
		target.httpClient = c
	}
}

// WithLogger sets the *logrus.Logger for the ClearBank API client.
func WithLogger(l *logrus.Logger) ClientOption {
	return func(target *client) {
		target.logger = l
	}
}

// WithBaseURL sets the base URL for the ClearBank API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(target *client) {
		target.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func NewClient(token string, options ...ClientOption) *client {
	c := &client{
		httpClient: http.DefaultClient,
		baseURL:    defaultBaseURL,
		token:      token,
		perPage:    100,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *client) newRequest(ctx context.Context, method, url string, body interface{}) (*request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var b []byte

	switch req.Method {
	case http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete:
		if body != nil {
			b, err = json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal payload: %w", err)
			}

			req.Body = io.NopCloser(bytes.NewReader(b))
		}
	}

	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.request.method":       req.Method,
			"http.request.url":          req.URL.String(),
			"http.request.body.content": string(b),
		}).Debug("quidax.client -> request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return NewRequest(req), nil
}

func (c *client) do(ctx context.Context, req *request) error {
	resp, err := c.httpClient.Do(req.req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(b))

	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.response.status_code":  resp.StatusCode,
			"http.response.body.content": string(b),
			"http.response.headers":      resp.Header,
		}).Debug("quidax.client -> response")
	}

	if !slices.Contains(req.expectedStatuses, resp.StatusCode) {
		unexpectedResponse := UnexpectedResponse{Status: resp.StatusCode, Body: string(b)}

		var errResponse ErrResponse
		if err := json.Unmarshal(b, &errResponse); err != nil {
			return unexpectedResponse
		}

		return errResponse
	}

	if req.decodeTo != nil {
		if err := json.NewDecoder(resp.Body).Decode(req.decodeTo); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
