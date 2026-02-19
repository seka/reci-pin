package request

type CreateRecipeRequest struct {
	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Memo   string  `json:"memo"`
	TagIDs []int64 `json:"tag_ids"`
}

type UpdateRecipeRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Memo string `json:"memo"`
}

type SearchRecipeRequest struct {
	Query  string  `json:"query"`
	TagIDs []int64 `json:"tag_ids"`
}

type AddTagsRequest struct {
	TagIDs []int64 `json:"tag_ids"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type CreateRecipeImageRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}
