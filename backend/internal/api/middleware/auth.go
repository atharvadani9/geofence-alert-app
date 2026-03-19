package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/atharvadani9/geofence-alert-app/internal/config"
	"github.com/atharvadani9/geofence-alert-app/internal/database"
	"github.com/atharvadani9/geofence-alert-app/internal/models"
	"github.com/atharvadani9/geofence-alert-app/pkg/jwt"
	"github.com/atharvadani9/geofence-alert-app/pkg/utils"
)

func Auth(db *database.DB, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Missing authorization header")
				utils.RespondError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Println("Invalid authorization header")
				utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			token := parts[1]
			claims, err := jwt.ValidateToken(token, cfg.JWT.Secret)
			if err != nil {
				log.Println(err)
				utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			userRepo := database.NewUserRepo(db)
			user, err := userRepo.GetUserByID(r.Context(), claims.UserID)
			if err != nil {
				log.Println(err)
				utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), models.UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
