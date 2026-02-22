package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/seka/reci-pin/backend/config"
)

func NewClient(cfg config.SearchEfunc NewClient(cfg config.SearchEngine) (*elasticsearch.TypedClient, error) {
	fmt.Printf("Initializing Elasticsearch client with addresses: %v\n", cfg.Addresses)

	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return client, nil
}
