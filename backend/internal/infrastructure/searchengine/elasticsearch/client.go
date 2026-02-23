package elasticsearch

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/seka/reci-pin/backend/config"
	"github.com/seka/reci-pin/backend/internal/infrastructure/searchengine"
)

func NewClient(cfg config.SearchEngine) (searchengine.SearchEngine, error) {
	fmt.Printf("Initializing Elasticsearch client with addresses: %v\n", cfg.Addresses)

	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}

	client, err := elasticsearch.NewTypedClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &searchEngineWrapper{client: client}, nil
}

type searchEngineWrapper struct {
	client *elasticsearch.TypedClient
}

func (w *searchEngineWrapper) Index(indexName string) searchengine.IndexService {
	return &indexServiceWrapper{service: w.client.Index(indexName)}
}

func (w *searchEngineWrapper) Delete(indexName string, id string) searchengine.DeleteService {
	return &deleteServiceWrapper{service: w.client.Delete(indexName, id)}
}

func (w *searchEngineWrapper) Search() searchengine.SearchService {
	return &searchServiceWrapper{service: w.client.Search()}
}

type indexServiceWrapper struct {
	service *index.Index
}

func (w *indexServiceWrapper) Id(id string) searchengine.IndexService {
	w.service = w.service.Id(id)
	return w
}

func (w *indexServiceWrapper) Request(doc interface{}) searchengine.IndexService {
	w.service = w.service.Request(doc)
	return w
}

func (w *indexServiceWrapper) Do(ctx context.Context) (*index.Response, error) {
	return w.service.Do(ctx)
}

type deleteServiceWrapper struct {
	service *delete.Delete
}

func (w *deleteServiceWrapper) Do(ctx context.Context) (*delete.Response, error) {
	return w.service.Do(ctx)
}

type searchServiceWrapper struct {
	service *search.Search
}

func (w *searchServiceWrapper) Index(indexName string) searchengine.SearchService {
	w.service = w.service.Index(indexName)
	return w
}

func (w *searchServiceWrapper) Request(req *search.Request) searchengine.SearchService {
	w.service = w.service.Request(req)
	return w
}

func (w *searchServiceWrapper) Do(ctx context.Context) (*search.Response, error) {
	return w.service.Do(ctx)
}
