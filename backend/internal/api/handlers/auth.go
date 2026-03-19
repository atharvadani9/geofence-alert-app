package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/atharvadani9/geofence-alert-app/internal/config"
	"github.com/atharvadani9/geofence-alert-app/internal/database"
	"github.com/atharvadani9/geofence-alert-app/internal/models"
	"github.com/atharvadani9/geofence-alert-app/pkg/jwt"
	"github.com/atharvadani9/geofence-alert-app/pkg/password"
	"github.com/atharvadani9/geofence-alert-app/pkg/utils"
	"github.com/go-playground/validator/v10"
)

func RegisterUser(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, "Failed to decode request body")
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		passwordHash, err := password.Hash(req.Password)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		userRepo := database.NewUserRepo(db)
		user, err := userRepo.CreateUser(r.Context(), req.Email, passwordHash, req.Name, req.Role)
		if err != nil {
			log.Println(err)
			if strings.Contains(err.Error(), "email already exists") {
				utils.RespondError(w, http.StatusConflict, "Email already exists")
				return
			}
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Email, user.Role, cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshTokenExpiry)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		resp := models.UserResponse{
			User:         *user,
			Token:        tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		utils.RespondJSON(w, http.StatusCreated, resp)
	}
}

func LoginUser(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, "Failed to decode request body")
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		userRepo := database.NewUserRepo(db)
		user, err := userRepo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		if err := password.Verify(user.PasswordHash, req.Password); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Email, user.Role, cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshTokenExpiry)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		resp := models.UserResponse{
			User:         *user,
			Token:        tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		utils.RespondJSON(w, http.StatusOK, resp)
	}
}

func RefreshToken(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.RefreshTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, "Failed to decode request body")
			return
		}

		if err := validate.Struct(req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		claims, err := jwt.ValidateToken(req.RefreshToken, cfg.JWT.Secret)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		userRepo := database.NewUserRepo(db)
		user, err := userRepo.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Email, user.Role, cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshTokenExpiry)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		resp := models.UserResponse{
			User:         *user,
			Token:        tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		utils.RespondJSON(w, http.StatusOK, resp)
	}
}

func GetCurrentUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user from context
		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		utils.RespondJSON(w, http.StatusOK, user)
	}
}
