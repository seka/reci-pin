package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/repository"
)

const indexName = "recipes"

type RecipeSearcher struct {
	client *elasticsearch.TypedClient
}

func NewRecipeSearcher(client *elasticsearch.TypedClient) repository.RecipeSearchRepository {
	return &RecipeSearcher{client: client}
}

type recipeDocument struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	Name      string  `json:"name"`
	Memo      string  `json:"memo"`
	TagIDs    []int64 `json:"tag_ids"`
	CreatedAt string  `json:"created_at"` // ISO8601 string
}

func (r *RecipeSearcher) Index(ctx context.Context, recipe *model.Recipe) error {
	tagIDs := make([]int64, len(recipe.Tags))
	for i, tag := range recipe.Tags {
		tagIDs[i] = tag.ID
	}

	doc := recipeDocument{
		ID:        recipe.ID,
		UserID:    recipe.UserID,
		Name:      recipe.Name,
		Memo:      recipe.Memo,
		TagIDs:    tagIDs,
		CreatedAt: recipe.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	_, err := r.client.Index(indexName).
		Id(strconv.FormatInt(recipe.ID, 10)).
		Request(doc).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to index recipe: %w", err)
	}

	return nil
}

func (r *RecipeSearcher) Delete(ctx context.Context, id int64) error {
	_, err := r.client.Delete(indexName, strconv.FormatInt(id, 10)).Do(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete recipe from index: %w", err)
	}
	return nil
}

func (r *RecipeSearcher) Search(ctx context.Context, criteria repository.SearchCriteria) ([]int64, int64, error) {
	mustQueries := []types.Query{
		{
			Term: map[string]types.TermQuery{
				"user_id": {Value: criteria.UserID},
			},
		},
	}

	if criteria.Keyword != "" {
		mustQueries = append(mustQueries, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  criteria.Keyword,
				Fields: []string{"name", "memo"},
			},
		})
	}

	if len(criteria.TagIDs) > 0 {
		// AND search for tags: document must have all specified tags
		for _, tagID := range criteria.TagIDs {
			mustQueries = append(mustQueries, types.Query{
				Term: map[string]types.TermQuery{
					"tag_ids": {Value: tagID},
				},
			})
		}
	}

	if criteria.Page <= 0 {
		criteria.Page = 1
	}
	if criteria.PageSize <= 0 {
		criteria.PageSize = 20
	}
	from := (criteria.Page - 1) * criteria.PageSize

	res, err := r.client.Search().
		Index(indexName).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must: mustQueries,
				},
			},
			From: &from,
			Size: &criteria.PageSize,
			Sort: []types.SortCombinations{
				types.SortOptions{
					SortOptions: map[string]types.FieldSort{
						"created_at": {Order: &sortorder.Desc},
					},
				},
			},
		}).
		Do(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed to search recipes: %w", err)
	}

	total := res.Hits.Total.Value
	ids := make([]int64, 0, len(res.Hits.Hits))
	for _, hit := range res.Hits.Hits {
		var doc recipeDocument
		if err := json.Unmarshal(hit.Source_, &doc); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal search hit: %w", err)
		}
		ids = append(ids, doc.ID)
	}

	return ids, total, nil
}
