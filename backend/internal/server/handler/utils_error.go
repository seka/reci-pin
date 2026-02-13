package handler

import (
	"net/http"

	"github.com/seka/reci-pin/backend/internal/server/handler/response"
)

func respondError(w http.ResponseWriter, status int, code string, details map[string][]response.ErrorDetail) {
	res := response.ErrorResponse{
		Error: response.ErrorContent{
			Code:    code,
			Details: details,
		},
	}
	respondJSON(w, status, res)
}
