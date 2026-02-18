package elasticsearch_test

import (
	"context"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
	es_repo "github.com/seka/reci-pin/backend/internal/infrastructure/database/elasticsearch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testIndexName = "recipes" // Using the same index name for simplicity, assume dev env
	esAddress     = "http://localhost:9200"
)

func setupTestClient(t *testing.T) *elasticsearch.TypedClient {
	cfg := elasticsearch.Config{
		Addresses: []string{esAddress},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	require.NoError(t, err)

	// Ping to ensure connection
	_, err = client.Info().Do(context.Background())
	if err != nil {
		t.Skipf("Skipping integration test: Elasticsearch not available at %s: %v", esAddress, err)
	}

	return client
}

func cleanupIndex(t *testing.T, client *elasticsearch.TypedClient) {
	_, err := client.Indices.Delete(testIndexName).IgnoreUnavailable(true).Do(context.Background())
	require.NoError(t, err)
	
	// Re-create is handled by migration/sync usually, but here we depend on auto-creation or existing index.
	// For "Index" test, mapping might be auto-created even if index doesn't exist, 
	// but strictly speaking we should probably ensure proper mapping.
	// For now, we'll just delete documents or let them be.
	// Actually, deleting the index is destructive for other data if any.
	// Better to just delete the specific test data.
	
	// Ideally we run this against a separate test index or clean up specific IDs.
}

func TestRecipeSearcher_Integration(t *testing.T) {
	client := setupTestClient(t)
	searcher := es_repo.NewRecipeSearcher(client)
	ctx := context.Background()

	// Test case data
	recipe := &model.Recipe{
		ID:        99999, // Use a high ID to avoid conflict
		UserID:    12345,
		Name:      "Integration Test Curry",
		Memo:      "Spicy and tasty test",
		CreatedAt: time.Now(),
		Tags: []model.Tag{
			{ID: 10, Name: "Spicy"},
			{ID: 11, Name: "Dinner"},
		},
	}

	t.Run("Index Recipe", func(t *testing.T) {
		err := searcher.Index(ctx, recipe)
		assert.NoError(t, err)

		// Refresh index to make document searchable immediately
		_, err = client.Indices.Refresh().Index(testIndexName).Do(ctx)
		require.NoError(t, err)
	})

	t.Run("Search Recipe by Keyword", func(t *testing.T) {
		// Eventual consistency might still apply slightly, but Refresh should handle it.
		
		ids, total, err := searcher.Search(ctx, repository.SearchCriteria{
			UserID:  recipe.UserID,
			Keyword: "Curry",
			Page:    1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		assert.Contains(t, ids, recipe.ID)

		// Test partial match
		ids, total, err = searcher.Search(ctx, repository.SearchCriteria{
			UserID:  recipe.UserID,
			Keyword: "Spicy", // Matches Memo or Tag? Search implementation only checks name/memo for keyword
			Page:    1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		assert.Contains(t, ids, recipe.ID)
	})

	t.Run("Search Recipe by Tag", func(t *testing.T) {
		ids, total, err := searcher.Search(ctx, repository.SearchCriteria{
			UserID: recipe.UserID,
			TagIDs: []int64{10}, // "Spicy"
			Page:   1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, total, int64(1))
		assert.Contains(t, ids, recipe.ID)
	})

	t.Run("Search Recipe - No Match", func(t *testing.T) {
		ids, total, err := searcher.Search(ctx, repository.SearchCriteria{
			UserID:  recipe.UserID,
			Keyword: "NonExistentThingy",
			Page:    1,
			PageSize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, ids)
	})

	t.Run("Delete Recipe", func(t *testing.T) {
		err := searcher.Delete(ctx, recipe.ID)
		assert.NoError(t, err)

		// Refresh
		_, err = client.Indices.Refresh().Index(testIndexName).Do(ctx)
		require.NoError(t, err)

		// Should not be found
		ids, _, err := searcher.Search(ctx, repository.SearchCriteria{
			UserID:  recipe.UserID,
			Keyword: "Curry",
		})
		assert.NoError(t, err)
		// We can't guarantee total is 0 if other tests added curries, but our ID should be gone
		assert.NotContains(t, ids, recipe.ID)
	})
}
