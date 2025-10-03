package quidax_test

import (
	"bytes"
	"context"
	_ "embed"
	"io"
	"net/http"
	"testing"

	"github.com/brokeyourbike/quidax-api-client-go"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//go:embed testdata/fees-one.json
var feesOne []byte

//go:embed testdata/fees-multi.json
var feesMulti []byte

func TestFetchWithdrawalFees_One(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(feesOne))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWithdrawalFees(context.TODO(), "usdt", "")
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
	assert.Len(t, got.GetFees(), 1)
}

func TestFetchWithdrawalFees_Multi(t *testing.T) {
	mockHttpClient := quidax.NewMockHttpClient(t)
	client := quidax.NewClient("token", quidax.WithHTTPClient(mockHttpClient))

	resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(feesMulti))}
	mockHttpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Once()

	got, err := client.FetchWithdrawalFees(context.TODO(), "usdt", "")
	require.NoError(t, err)
	assert.Equal(t, "success", got.Status)
	assert.Len(t, got.GetFees(), 7)
}
