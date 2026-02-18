package scenario

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	Register("search_recipes", SearchRecipes)
}

func SearchRecipes(ctx context.Context, client *http.Client, baseURL string) error {
	// 1. Login
	// Using seeded user "john@example.com"
	loginData := map[string]string{
		"email":    "john@example.com",
		"password": "password",
	}
	b, err := json.Marshal(loginData)
	if err != nil {
		return err
	}

	resp, err := client.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login status %d", resp.StatusCode)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return fmt.Errorf("decode login: %w", err)
	}

	if loginResp.Token == "" {
		return fmt.Errorf("empty token")
	}

	// 2. Search
	searchData := map[string]string{
		"query": "Fish", // Seed data has "Delicious Fish"
	}
	sb, err := json.Marshal(searchData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/recipes/search", bytes.NewBuffer(sb))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	sResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer sResp.Body.Close()

	if sResp.StatusCode != http.StatusOK {
		return fmt.Errorf("search status %d", sResp.StatusCode)
	}

	// Optional: decode and verify results count?
	// For stress test, success status is enough.

	return nil
}
