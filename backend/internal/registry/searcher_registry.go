package registry

import (
	"github.com/seka/reci-pin/backend/internal/domain/searcher"
	infraSearchEngine "github.com/seka/reci-pin/backend/internal/infrastructure/searchengine"
	"github.com/seka/reci-pin/backend/internal/infrastructure/searchengine/elasticsearch"
)

// Searcher defines the interface for creating searchers
type Searcher interface {
	NewRecipeSearchRepository() searcher.RecipeSearcher
}

// searcherRegistry implements the Searcher interface
type searcherRegistry struct {
	searchEngine infraSearchEngine.SearchEngine
}

// NewSearcher creates a new Searcher registry
func NewSearcher(searchEngine infraSearchEngine.SearchEngine) Searcher {
	return &searcherRegistry{
		searchEngine: searchEngine,
	}
}

func (r *searcherRegistry) NewRecipeSearchRepository() searcher.RecipeSearcher {
	return elasticsearch.NewRecipeSearcher(r.searchEngine)
}
