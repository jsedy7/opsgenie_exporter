package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// OpsgenieClient represents a client for the Opsgenie API.
type OpsgenieClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// ListUsersResponse wraps the response structure for the ListUsers API call.
type ListUsersResponse struct {
	Data   []User `json:"data"` // List of users returned by the API.
	Paging struct {
		Next string `json:"next"` // URL for the next page of users, if any.
	} `json:"paging"`
}

// User represents a user object in Opsgenie.
type User struct {
	Username string `json:"username"` // Opsgenie username of the user.
	Blocked  bool   `json:"blocked"`  // Indicates if the user is blocked.
	Verified bool   `json:"verified"` // Indicates if the user is verified.
}

// ListTeamsResponse wraps the response structure for the ListTeams API call.
type ListTeamsResponse struct {
	Data   []Team `json:"data"` // List of teams returned by the API.
	Paging struct {
		Next string `json:"next"` // URL for the next page of teams, if any.
	} `json:"paging"`
}

// Team represents a team object in Opsgenie.
type Team struct {
	Name string `json:"name"` // Name of the team.
}

// Integration represents an integration object in Opsgenie.
type Integration struct {
	ID      string `json:"id"`      // Unique identifier of the integration.
	Name    string `json:"name"`    // Name of the integration.
	Enabled bool   `json:"enabled"` // Indicates if the integration is enabled.
	Type    string `json:"type"`    // Type of the integration.
}

// ListIntegrationsResponse wraps the response structure for the ListIntegrations API call.
type ListIntegrationsResponse struct {
	Data []Integration `json:"data"` // List of integrations returned by the API.
}

// Heartbeat represents a heartbeat object in Opsgenie.
type Heartbeat struct {
	Name         string `json:"name"`         // Name of the heartbeat.
	Description  string `json:"description"`  // Description of the heartbeat.
	Interval     int    `json:"interval"`     // Interval of the heartbeat.
	IntervalUnit string `json:"intervalUnit"` // Unit of the interval (e.g., minutes, hours).
	Enabled      bool   `json:"enabled"`      // Indicates if the heartbeat is enabled.
	OwnerTeam    struct {
		ID   string `json:"id"`   // ID of the team owning the heartbeat.
		Name string `json:"name"` // Name of the team owning the heartbeat.
	} `json:"ownerTeam"`
	AlertMessage  string   `json:"alertMessage"`  // Alert message for the heartbeat.
	AlertTags     []string `json:"alertTags"`     // Tags associated with the heartbeat alert.
	AlertPriority string   `json:"alertPriority"` // Priority of the heartbeat alert.
}

// ListHeartbeatsResponse wraps the response structure for the ListHeartbeats API call.
type ListHeartbeatsResponse struct {
	Data struct {
		Heartbeats []Heartbeat `json:"heartbeats"` // List of heartbeats returned by the API.
	} `json:"data"`
}

// HeartbeatDetail represents detailed information about a heartbeat in Opsgenie.
type HeartbeatDetail struct {
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	Expired      bool   `json:"expired"`
	Interval     int    `json:"interval"`
	IntervalUnit string `json:"intervalUnit"`
	OwnerTeam    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"ownerTeam"`
	AlertMessage  string   `json:"alertMessage"`
	AlertTags     []string `json:"alertTags"`
	AlertPriority string   `json:"alertPriority"`
}

// AccountInfoResponse wraps the response structure for fetching account information.
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

// NewOpsgenieClient initializes a new Opsgenie API client with the provided API key.
func NewOpsgenieClient(apiKey string) *OpsgenieClient {
	return &OpsgenieClient{
		APIKey:     apiKey,
		BaseURL:    "https://api.eu.opsgenie.com/v2/",
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// callAPI makes a GET request to the specified Opsgenie API endpoint and returns the response body.
// It automatically adds the required authorization header.
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

// ListUsers fetches a list of users from Opsgenie and handles pagination.
func (c *OpsgenieClient) ListUsers() ([]User, error) {
	var allUsers []User
	next := "users?limit=100" // Initial page

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

		// Process pagination if the response includes information about the next page.
		if response.Paging.Next != "" {
			next = strings.TrimPrefix(response.Paging.Next, c.BaseURL)
		} else {
			next = "" // End of pagination
		}
	}
	return allUsers, nil
}

// ListTeams fetches a list of teams from Opsgenie and handles pagination.
func (c *OpsgenieClient) ListTeams() ([]Team, error) {
	var allTeams []Team
	next := "teams?limit=100" // Initial page

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

		// Process pagination if the response includes information about the next page.
		if response.Paging.Next != "" {
			next = strings.TrimPrefix(response.Paging.Next, c.BaseURL)
		} else {
			next = "" // End of pagination
		}
	}

	return allTeams, nil
}

// GetAccountInfo fetches information about the Opsgenie account.
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

// ListIntegrations fetches a list of integrations from Opsgenie.
func (c *OpsgenieClient) ListIntegrations() ([]Integration, error) {
	body, err := c.callAPI("integrations")
	if err != nil {
		return nil, err
	}

	var response ListIntegrationsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

// ListHeartbeats fetches a list of heartbeats from Opsgenie.
func (c *OpsgenieClient) ListHeartbeats() ([]Heartbeat, error) {
	body, err := c.callAPI("heartbeats")
	if err != nil {
		return nil, err
	}

	var response ListHeartbeatsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Data.Heartbeats, nil
}

// GetHeartbeatDetail fetches detailed information about a specific heartbeat from Opsgenie.
func (c *OpsgenieClient) GetHeartbeatDetail(name string) (*HeartbeatDetail, error) {
	endpoint := fmt.Sprintf("heartbeats/%s", name)
	body, err := c.callAPI(endpoint)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data HeartbeatDetail `json:"data"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
