package scenario

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	Register("auth_flow", AuthFlow)
}

func AuthFlow(ctx context.Context, client *http.Client, baseURL string) error {
	// 1. Signup
	// Use random email to avoid conflict? Or reuse?
	// If allow repeats, reuse. But Signup usually fails on duplicate.
	// Use unique email per iteration?
	email := fmt.Sprintf("stress_%d_%d@example.com", time.Now().UnixNano(), rand.Int())
	password := "password"

	signupData := map[string]string{"name": "Stress User", "email": email, "password": password}
	if err := postJSON(client, baseURL+"/signup", signupData); err != nil {
		return err
	}

	// 2. Login
	loginData := map[string]string{"email": email, "password": password}
	if err := postJSON(client, baseURL+"/login", loginData); err != nil {
		return err
	}

	// 3. Get Profile (Admin area? Or just verify login?)
	// Reci-pin might not have /me yet.
	// Doing nothing implies successfully logged in.
	return nil
}
