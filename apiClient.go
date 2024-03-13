package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type OpsgenieClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

type ListUsersResponse struct {
	Data   []User `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type User struct {
	Username string `json:"username"`
	Blocked  bool   `json:"blocked"`
	Verified bool   `json:"verified"`
}

type ListTeamsResponse struct {
	Data   []Team `json:"data"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

type Team struct {
	Name string `json:"name"`
}

type AccountInfoResponse struct {
	Data struct {
		Name      string `json:"name"`
		UserCount int    `json:"userCount"`
		Plan      struct {
			MaxUserCount int    `json:"maxUserCount"`
			Name         string `json:"name"`
			IsYearly     bool   `json:"isYearly"`
		} `json:"plan"`
	} `json:"data"`
}

func NewOpsgenieClient(apiKey string) *OpsgenieClient {
	return &OpsgenieClient{
		APIKey:     apiKey,
		BaseURL:    "https://api.eu.opsgenie.com/v2/",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *OpsgenieClient) callAPI(endpoint string) ([]byte, error) {
	requestURL := c.BaseURL + endpoint

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "GenieKey "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *OpsgenieClient) ListUsers() ([]User, error) {
	var allUsers []User
	next := "users?limit=100" // Počáteční stránka

	for next != "" {
		body, err := c.callAPI(next)
		if err != nil {
			return nil, err
		}

		var response ListUsersResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allUsers = append(allUsers, response.Data...)

		// Zpracování stránkování. Příklad předpokládá, že odpověď obsahuje informace o další stránce.
		// Musíte upravit logiku podle struktury odpovědi API.
		if response.Paging.Next != "" {
			// Příklad získání další stránky z URL. Může vyžadovat úpravu.
			next = strings.TrimPrefix(response.Paging.Next, c.BaseURL)
		} else {
			next = "" // Konec stránkování
		}
	}
	return allUsers, nil
}

func (c *OpsgenieClient) ListTeams() ([]Team, error) {
	var allTeams []Team
	next := "teams?limit=100" // Počáteční stránka

	for next != "" {
		body, err := c.callAPI(next)
		if err != nil {
			return nil, err
		}

		var response ListTeamsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, err
		}

		allTeams = append(allTeams, response.Data...)

		// Zpracování stránkování. Předpokládá, že odpověď obsahuje informace o další stránce.
		if response.Paging.Next != "" {
			// Získání další stránky z URL. Může vyžadovat úpravu.
			next = strings.TrimPrefix(response.Paging.Next, c.BaseURL)
		} else {
			next = "" // Konec stránkování
		}
	}

	return allTeams, nil
}

func (c *OpsgenieClient) GetAccountInfo() (*AccountInfoResponse, error) {
	body, err := c.callAPI("account")
	if err != nil {
		return nil, err
	}

	var response AccountInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
