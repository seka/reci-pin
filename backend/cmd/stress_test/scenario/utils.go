package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func postJSON(client *http.Client, url string, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
