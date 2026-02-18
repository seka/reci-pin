package elasticsearch

import (
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewClient() (*elasticsearch.TypedClient, error) {
	esAddress := os.Getenv("ELASTICSEARCH_ADDRESS")
	if esAddress == "" {
		esAddress = "http://localhost:9200"
	}
	fmt.Printf("Initializing Elasticsearch client with address: %s\n", esAddress)

	cfg := elasticsearch.Config{
		Addresses: []string{esAddress},
	}

	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return client, nil
}
