package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-service/internal/model"
	"project-service/internal/service/customServiceError"
	"time"
)

func GetUserByToken(token string, api, port string) (*model.User, error) {
	client := http.Client{Timeout: time.Second * 30}
	url := fmt.Sprintf("http://%s:%s/api/user/get", api, port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		var errResp struct {
			Status bool   `json:"status"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, customServiceError.ErrTokenNotFound
		}
		return nil, customServiceError.ErrTokenNotFound
	}
	if resp.StatusCode == http.StatusUnauthorized {
		var errResp struct {
			Status bool   `json:"status"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, customServiceError.ErrTokenExpired
		}
		return nil, customServiceError.ErrTokenExpired
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id string, api, port string) (*model.User, error) {
	client := http.Client{Timeout: time.Second * 30}

	url := fmt.Sprintf("http://%s:%s/api/outside/user/get/%s", api, port, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
