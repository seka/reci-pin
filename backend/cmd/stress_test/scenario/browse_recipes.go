package scenario

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
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
	resp2, err := client.Get(fmt.Sprintf("%s/recipes/%d", baseURL, id))
	if err != nil {
		return err
	}
	defer func() { _ = resp2.Body.Close() }()
	// Ignore status 404 for random IDs

	return nil
}
