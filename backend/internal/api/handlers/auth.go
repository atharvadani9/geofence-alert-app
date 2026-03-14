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

var validate = validator.New()

func RegisterUser(db *database.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.UserRegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, "Failed to decode request body")
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Hash password
		passwordHash, err := password.Hash(req.Password)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		// Create user
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

		// Generate tokens
		tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Email, user.Role, cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshTokenExpiry)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		// Return response
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
		var req models.UserLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, "Failed to decode request body")
			return
		}

		// Validate request
		if err := validate.Struct(req); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get user
		userRepo := database.NewUserRepo(db)
		user, err := userRepo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Verify password
		if err := password.Verify(user.PasswordHash, req.Password); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Generate tokens
		tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Email, user.Role, cfg.JWT.Secret, cfg.JWT.Expiry, cfg.JWT.RefreshTokenExpiry)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		// Return response
		resp := models.UserResponse{
			User:         *user,
			Token:        tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		}
		utils.RespondJSON(w, http.StatusOK, resp)
	}
}
