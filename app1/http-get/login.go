package http_get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goDevOpCloud/utils"
	"io"
	"net/http"
)

type LoginRequest struct {
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func doLoginRequest(client http.Client, requestURL, password string) (string, error) {
	loginRequest := LoginRequest{
		Password: password,
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", err)
	}

	resp, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", fmt.Errorf("HTTP Post error: %s", err)
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("ReadAll error: %s", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid output (HTTP Code %d): %s\n", resp.StatusCode, string(resBody))
	}

	var loginResponse LoginResponse

	if !json.Valid(resBody) {
		return "", utils.RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("No valid JSON returned"),
		}
	}

	err = json.Unmarshal(resBody, &loginResponse)
	if err != nil {
		return "", utils.RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      fmt.Sprintf("Unmarshal error: %s", err),
		}
	}

	if loginResponse.Token == "" {
		return "", utils.RequestError{
			HTTPCode: resp.StatusCode,
			Body:     string(resBody),
			Err:      "Empty token replied",
		}
	}

	return loginResponse.Token, nil
}
