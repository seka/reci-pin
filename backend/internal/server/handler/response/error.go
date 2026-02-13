package response

type ErrorDetail struct {
	Code   string                 `json:"code"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type ErrorContent struct {
	Code    string                   `json:"code"`
	Details map[string][]ErrorDetail `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}
