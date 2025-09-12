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

//go:embed testdata/accounts-fetch-me-ok.json
var accountsFetchMeOK []byte

//go:embed testdata/accounts-fetch-all-ok.json
var accountsFetchAllOK []byte

func TestFetchParentAccount_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountsFetchMeOK))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchParentAccount(context.TODO())
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
}

func TestFetchAccount_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountsFetchMeOK))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchAccount(context.TODO(), uuid.New())
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
}

func TestFetchAccounts_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountsFetchAllOK))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchAccounts(context.TODO(), 0)
	require.NoError(t, err)
	assert.Len(t, got.Data, 1)
}

func TestCreateAccount_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewReader(accountsFetchMeOK))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.CreateAccount(context.TODO(), quidax.CreateAccountPayload{})
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
}

func TestUpdateAccount_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(accountsFetchMeOK))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.UpdateAccount(context.TODO(), uuid.New(), quidax.UpdateAccountPayload{})
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
}
