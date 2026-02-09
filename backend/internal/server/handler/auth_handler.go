package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/server/handler/request"
	"github.com/seka/reci-pin/backend/internal/server/handler/response"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
)

type AuthHandler struct {
	signupUseCase         auth.SignupUseCase
	loginUseCase          auth.LoginUseCase
	generateTokenUseCase  auth.GenerateTokenUseCase
	getUserUseCase        auth.GetUserUseCase
	verifyEmailUseCase    auth.VerifyEmailUseCase
	withdrawUseCase       auth.WithdrawUseCase
	changePasswordUseCase auth.ChangePasswordUseCase
}

func NewAuthHandler(
	signupUseCase auth.SignupUseCase,
	loginUseCase auth.LoginUseCase,
	generateTokenUseCase auth.GenerateTokenUseCase,
	getUserUseCase auth.GetUserUseCase,
	verifyEmailUseCase auth.VerifyEmailUseCase,
	withdrawUseCase auth.WithdrawUseCase,
	changePasswordUseCase auth.ChangePasswordUseCase,
) *AuthHandler {
	return &AuthHandler{
		signupUseCase:         signupUseCase,
		loginUseCase:          loginUseCase,
		generateTokenUseCase:  generateTokenUseCase,
		getUserUseCase:        getUserUseCase,
		verifyEmailUseCase:    verifyEmailUseCase,
		withdrawUseCase:       withdrawUseCase,
		changePasswordUseCase: changePasswordUseCase,
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

	respondJSON(w, http.StatusCreated, res)
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

	respondJSON(w, http.StatusOK, res)
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

	respondJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.withdrawUseCase.Execute(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req request.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := auth.ChangePasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	if err := h.changePasswordUseCase.Execute(r.Context(), userID, input); err != nil {
		// エラーの内容によってステータスコードを変えるべきだが、簡単のため400にする
		// 実際には internal error と bad request を区別すべき
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := response.MessageResponse{
		Message: "Password changed successfully",
	}

	respondJSON(w, http.StatusOK, res)
}
