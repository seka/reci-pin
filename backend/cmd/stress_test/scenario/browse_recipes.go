package scenario

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	Register("browse_recipes", BrowseRecipes)
}

func BrowseRecipes(ctx context.Context, client *http.Client, baseURL string) error {
	// 1. Get Recipes List
	resp, err := client.Get(baseURL + "/recipes")
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != 200 {
		return fmt.Errorf("list recipes status: %d", resp.StatusCode)
	}

	// 2. Simulate getting a detail (mock ID 1)
	// In real stress test, parse IDs from list.
	// For now just random 1-5
	id := rand.Intn(5) + 1
	recipePath, err := url.JoinPath("recipes", strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("creating recipe detail path: %w", err)
	}
	detailURL, err := url.JoinPath(baseURL, recipePath)
	if err != nil {
		return fmt.Errorf("creating detail URL: %w", err)
	}
	resp2, err := client.Get(detailURL)
	if err != nil {
		return err
	}
	defer func() { _ = resp2.Body.Close() }()
	// Ignore status 404 for random IDs

	return nil
}
