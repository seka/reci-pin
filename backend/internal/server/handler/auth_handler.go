package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/domain/validation"
	"github.com/seka/reci-pin/backend/internal/server/handler/request"
	"github.com/seka/reci-pin/backend/internal/server/handler/response"
	"github.com/seka/reci-pin/backend/internal/server/middleware"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
)

type AuthHandler struct {
	signupUseCase               auth.SignupUseCase
	loginUseCase                auth.LoginUseCase
	generateTokenUseCase        auth.GenerateTokenUseCase
	getUserUseCase              auth.GetUserUseCase
	verifyEmailUseCase          auth.VerifyEmailUseCase
	withdrawUseCase             auth.WithdrawUseCase
	changePasswordUseCase       auth.ChangePasswordUseCase
	requestPasswordResetUseCase auth.RequestPasswordResetUseCase
	resetPasswordUseCase        auth.ResetPasswordUseCase
	refreshTokenUseCase         auth.RefreshTokenUseCase
	logoutUseCase               auth.LogoutUseCase
}

func NewAuthHandler(
	signupUseCase auth.SignupUseCase,
	loginUseCase auth.LoginUseCase,
	generateTokenUseCase auth.GenerateTokenUseCase,
	getUserUseCase auth.GetUserUseCase,
	verifyEmailUseCase auth.VerifyEmailUseCase,
	withdrawUseCase auth.WithdrawUseCase,
	changePasswordUseCase auth.ChangePasswordUseCase,
	requestPasswordResetUseCase auth.RequestPasswordResetUseCase,
	resetPasswordUseCase auth.ResetPasswordUseCase,
	refreshTokenUseCase auth.RefreshTokenUseCase,
	logoutUseCase auth.LogoutUseCase,
) *AuthHandler {
	return &AuthHandler{
		signupUseCase:               signupUseCase,
		loginUseCase:                loginUseCase,
		generateTokenUseCase:        generateTokenUseCase,
		getUserUseCase:              getUserUseCase,
		verifyEmailUseCase:          verifyEmailUseCase,
		withdrawUseCase:             withdrawUseCase,
		changePasswordUseCase:       changePasswordUseCase,
		requestPasswordResetUseCase: requestPasswordResetUseCase,
		resetPasswordUseCase:        resetPasswordUseCase,
		refreshTokenUseCase:         refreshTokenUseCase,
		logoutUseCase:               logoutUseCase,
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
	log.Println("Handling Signup request")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Signup handler: %v", r)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}()

	var req request.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
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
		log.Printf("SignupUseCase error: %v", err)
		var validationErrors validation.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := make(map[string][]response.ErrorDetail)
			for _, ve := range validationErrors {
				details[ve.Field] = append(details[ve.Field], response.ErrorDetail{
					Code:   ve.Code,
					Params: ve.Params,
				})
			}
			respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", details)
			return
		}
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

	// 認証トークン生成 (Access + Refresh)
	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr

	tokenResult, err := h.generateTokenUseCase.Execute(userID, userAgent, ipAddress)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	// HttpOnly Cookie 設定
	h.setAuthCookies(w, tokenResult)

	userResp := toUserResponse(user)
	// UserモデルにEmailがないため、リクエストの値を使用
	userResp.Email = req.Email

	res := response.AuthResponse{
		Token: "", // トークンはCookieに隠蔽するため、ボディには含めない
		User:  userResp,
	}

	respondJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	userAgent := r.UserAgent()
	ipAddress := r.RemoteAddr

	tokenResult, err := h.refreshTokenUseCase.Execute(cookie.Value, userAgent, ipAddress)
	if err != nil {
		log.Printf("RefreshTokenUseCase error: %v", err)
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	// 新しいトークンペアを Cookie に設定 (Rotation)
	h.setAuthCookies(w, tokenResult)

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, tokens *auth.TokenResult) {
	// Access Token
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  tokens.AccessTokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Refresh Token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		Expires:  tokens.RefreshTokenExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
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
	// DB側のリフレッシュトークンを無効化（もしあれば）
	if cookie, err := r.Cookie("refresh_token"); err == nil {
		if err := h.logoutUseCase.Execute(cookie.Value); err != nil {
			log.Printf("LogoutUseCase error: %v", err)
			// 失敗してもクライアント側は消したいので続行
		}
	}

	// Cookieを削除（有効期限を過去に設定）
	h.clearAuthCookies(w)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	for _, name := range []string{"auth_token", "refresh_token"} {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
	}
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

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var req request.RequestPasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.requestPasswordResetUseCase.Execute(r.Context(), req.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := response.MessageResponse{
		Message: "メールアドレスが登録されている場合、パスワード再設定用のリンクをお送りします。",
	}

	respondJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req request.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := auth.ResetPasswordInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}

	if err := h.resetPasswordUseCase.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := response.MessageResponse{
		Message: "パスワードを再設定しました。ログイン画面へ移動します...",
	}

	respondJSON(w, http.StatusOK, res)
}
