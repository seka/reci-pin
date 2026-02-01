package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler/request"
	"github.com/seka/reci-pin/backend/internal/server/handler/response"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
)

type AuthHandler struct {
	signupUseCase        auth.SignupUseCase
	loginUseCase         auth.LoginUseCase
	generateTokenUseCase auth.GenerateTokenUseCase
	getUserUseCase       auth.GetUserUseCase
	verifyEmailUseCase   auth.VerifyEmailUseCase
}

func NewAuthHandler(
	signupUseCase auth.SignupUseCase,
	loginUseCase auth.LoginUseCase,
	generateTokenUseCase auth.GenerateTokenUseCase,
	getUserUseCase auth.GetUserUseCase,
	verifyEmailUseCase auth.VerifyEmailUseCase,
) *AuthHandler {
	return &AuthHandler{
		signupUseCase:        signupUseCase,
		loginUseCase:         loginUseCase,
		generateTokenUseCase: generateTokenUseCase,
		getUserUseCase:       getUserUseCase,
		verifyEmailUseCase:   verifyEmailUseCase,
	}
}

// Converters

func toUserResponse(user *model.User) *response.UserResponse {
	if user == nil {
		return nil
	}
	return &response.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		// EmailはUser(Profile)モデルに含まれないため、呼び出し元で設定する必要がある
	}
}

// Handlers

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req request.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := auth.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// 登録完了メッセージのみ返却（トークンはメール認証後に発行）
	_, err := h.signupUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := response.MessageResponse{
		Message: "Verification email sent. Please check your inbox.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	userID, err := h.loginUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// ユーザープロファイル取得
	user, err := h.getUserUseCase.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	// 認証トークン生成
	token, err := h.generateTokenUseCase.Execute(userID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	userResp := toUserResponse(user)
	// UserモデルにEmailがないため、リクエストの値を使用
	userResp.Email = req.Email

	res := response.AuthResponse{
		Token: token,
		User:  userResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req request.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.verifyEmailUseCase.Execute(r.Context(), req.Token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := response.MessageResponse{
		Message: "Email verified successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
