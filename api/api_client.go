package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIResponse struct {
	Status string `json:"status"`
	Detail string `json:"detail"`
}

func CheckMessageSent(apiURL string, messageID string) (*APIResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/sent/%s", apiURL, messageID))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result APIResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckMessageArrived(apiURL string, messageID string) (*APIResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/arrived/%s", apiURL, messageID))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result APIResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CheckMessageProcessed(apiURL string, messageID string) (*APIResponse, error) {
	response, err := http.Get(fmt.Sprintf("%s/processed/%s", apiURL, messageID))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result APIResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
