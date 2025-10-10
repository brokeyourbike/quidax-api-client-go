package quidax_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/quidax-api-client-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/swap-quote-ok.json
var swapQuoteOk []byte

func TestQuote_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewReader(swapQuoteOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.Quote(context.TODO(), uuid.New(), quidax.QuotePayload{})
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
	assert.Equal(t, 0.01, got.Data.GetFromAmount())
	assert.Equal(t, 0.00034847, got.Data.GetToAmount())
}

func TestConfirmQuote_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(swapQuoteOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	err := client.ConfirmQuote(context.TODO(), uuid.New(), uuid.New())
	require.NoError(t, err)
}
