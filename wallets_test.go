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

//go:embed testdata/wallets-fetch-btc-ok.json
var walletsFetchBtcOk []byte

//go:embed testdata/wallets-fetch-all-ok.json
var walletsFetchAllOk []byte

//go:embed testdata/wallets-address-fetch-all-ok.json
var walletsAddressFetchAllOk []byte

//go:embed testdata/wallets-address-fetch-btc-ok.json
var walletsAddressFetchBtcOk []byte

func TestFetchWallet_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(walletsFetchBtcOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWallet(context.TODO(), uuid.New(), "btc")
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
	assert.Equal(t, 10.00, got.Data.GetBalance())
	assert.Equal(t, 5.00, got.Data.GetLocked())
	assert.Equal(t, 1.00, got.Data.GetStaked())
}

func TestFetchWallets_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(walletsFetchAllOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWallets(context.TODO(), uuid.New())
	require.NoError(t, err)
	assert.Len(t, got.Data, 1)
}

func TestFetchWalletAddress_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(walletsAddressFetchBtcOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWalletAddress(context.TODO(), uuid.New(), "btc")
	require.NoError(t, err)
	assert.Equal(t, "dummyaddress", got.Data.Address)
}

func TestFetchWalletAddresses_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(walletsAddressFetchAllOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWalletAddresses(context.TODO(), uuid.New(), "btc")
	require.NoError(t, err)
	assert.Len(t, got.Data, 1)
}

func TestRequestWalletAddress_Success(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewReader(walletsAddressFetchBtcOk))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.RequestWalletAddress(context.TODO(), uuid.New(), "btc", "bep20")
	require.NoError(t, err)
	assert.Equal(t, "dummyaddress", got.Data.Address)
}
