package quidax

import "fmt"

type UnexpectedResponse struct {
	Status int
	Body   string
}

func (r UnexpectedResponse) Error() string {
	return fmt.Sprintf("Unexpected response from API. Status: %d Body: %s", r.Status, r.Body)
}

type ErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"data"`
}

func (e ErrResponse) Error() string {
	return fmt.Sprintf("Error during API call. Status: %s Message: %s", e.Status, e.Message)
}
