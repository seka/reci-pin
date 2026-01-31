package handler

import (
	"encoding/json"
	"net/http"

	"github.com/seka/reci-pin/backend/internal/domain/model"
	"github.com/seka/reci-pin/backend/internal/usecase/auth"
)

type AuthHandler struct {
	signupUseCase        *auth.SignupUseCase
	loginUseCase         *auth.LoginUseCase
	generateTokenUseCase *auth.GenerateTokenUseCase
	getUserUseCase       *auth.GetUserUseCase
	verifyEmailUseCase   *auth.VerifyEmailUseCase
}

func NewAuthHandler(
	signupUseCase *auth.SignupUseCase,
	loginUseCase *auth.LoginUseCase,
	generateTokenUseCase *auth.GenerateTokenUseCase,
	getUserUseCase *auth.GetUserUseCase,
	verifyEmailUseCase *auth.VerifyEmailUseCase,
) *AuthHandler {
	return &AuthHandler{
		signupUseCase:        signupUseCase,
		loginUseCase:         loginUseCase,
		generateTokenUseCase: generateTokenUseCase,
		getUserUseCase:       getUserUseCase,
		verifyEmailUseCase:   verifyEmailUseCase,
	}
}

// Request/Response structures

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type AuthResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// Converters

func toUserResponse(user *model.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:   user.ID,
		Name: user.Name,
		// Email must be filled by caller as it's not in User model anymore
	}
}

// Handlers

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	input := auth.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	_, err := h.signupUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := MessageResponse{
		Message: "Verification email sent. Please check your inbox.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
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

	// Get user
	user, err := h.getUserUseCase.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := h.generateTokenUseCase.Execute(userID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	userResp := toUserResponse(user)
	userResp.Email = req.Email // Set email from request

	response := AuthResponse{
		Token: token,
		User:  userResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.verifyEmailUseCase.Execute(r.Context(), req.Token); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := MessageResponse{
		Message: "Email verified successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
