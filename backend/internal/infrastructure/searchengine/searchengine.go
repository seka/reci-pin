package searchengine

import (
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/delete"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/index"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
)

// SearchEngine defines the interface for the search engine client,
// equivalent to the elasticsearch.TypedClient.
type SearchEngine interface {
	Index(index string) IndexService
	Delete(indexName string, id string) DeleteService
	Search() SearchService
}

// IndexService defines the interface for the index operation.
type IndexService interface {
	Id(id string) IndexService
	Request(doc any) IndexService
	Do(ctx context.Context) (*index.Response, error)
}

// DeleteService defines the interface for the delete operation.
type DeleteService interface {
	Do(ctx context.Context) (*delete.Response, error)
}

// SearchService defines the interface for the search operation.
type SearchService interface {
	Index(index string) SearchService
	Request(req *search.Request) SearchService
	Do(ctx context.Context) (*search.Response, error)
}
