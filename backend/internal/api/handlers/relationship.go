package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/atharvadani9/geofence-alert-app/internal/database"
	"github.com/atharvadani9/geofence-alert-app/internal/models"
	"github.com/atharvadani9/geofence-alert-app/pkg/utils"
	"github.com/go-playground/validator/v10"
)

func InviteUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.InviteRequest
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

		// Get current user from context
		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		// Check if user is caregiver
		if user.Role != "caregiver" {
			log.Println("User is not a caregiver")
			utils.RespondError(w, http.StatusForbidden, "Only caregivers can invite users")
			return
		}

		userRepo := database.NewUserRepo(db)
		trackedUser, err := userRepo.GetUserByEmail(r.Context(), req.TrackedUserEmail)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusNotFound, "Tracked user not found")
			return
		}

		// Check if relationship already exists
		relationshipRepo := database.NewRelationshipRepo(db)
		exists, err := relationshipRepo.CheckRelationshipExists(r.Context(), user.ID, trackedUser.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to check relationship")
			return
		}
		if exists {
			log.Println("Relationship already exists")
			utils.RespondError(w, http.StatusConflict, "Relationship already exists")
			return
		}

		relationship, err := relationshipRepo.CreateInvite(r.Context(), user.ID, trackedUser.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create invite")
			return
		}

		utils.RespondJSON(w, http.StatusCreated, relationship)
	}
}

func UpdateInviteStatus(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.UpdateInviteRequest
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

		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		relationshipRepo := database.NewRelationshipRepo(db)
		relationship, err := relationshipRepo.GetRelationship(r.Context(), req.RelationshipID, user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusNotFound, "Relationship not found")
			return
		}
		if relationship.Status != "pending" {
			log.Println("Relationship is not pending")
			utils.RespondError(w, http.StatusConflict, "Relationship is not pending")
			return
		}

		if req.Status == "accepted" {
			if err := relationshipRepo.AcceptInvite(r.Context(), req.RelationshipID); err != nil {
				log.Println(err)
				utils.RespondError(w, http.StatusInternalServerError, "Failed to accept invite")
				return
			}
		} else {
			if err := relationshipRepo.RejectInvite(r.Context(), req.RelationshipID); err != nil {
				log.Println(err)
				utils.RespondError(w, http.StatusInternalServerError, "Failed to reject invite")
				return
			}
		}

		utils.RespondJSON(w, http.StatusOK, "Invite updated")
	}
}

func DeleteRelationship(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.DeleteRelationshipRequest
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

		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		relationshipRepo := database.NewRelationshipRepo(db)
		relationship, err := relationshipRepo.GetRelationship(r.Context(), req.RelationshipID, user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusNotFound, "Relationship not found")
			return
		}

		if relationship.Status != "accepted" {
			log.Println("Relationship is not accepted")
			utils.RespondError(w, http.StatusConflict, "Relationship is not accepted")
			return
		}

		if err := relationshipRepo.DeleteRelationship(r.Context(), req.RelationshipID); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to delete relationship")
			return
		}

		utils.RespondJSON(w, http.StatusOK, "Relationship deleted")
	}
}

func GetCaregivers(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		relationshipRepo := database.NewRelationshipRepo(db)
		caregivers, err := relationshipRepo.GetCaregivers(r.Context(), user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get caregivers")
			return
		}

		utils.RespondJSON(w, http.StatusOK, caregivers)
	}
}

func GetTrackedUsers(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(models.UserContextKey).(*models.User)
		if !ok {
			log.Println("Failed to get user from context")
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get user from context")
			return
		}

		relationshipRepo := database.NewRelationshipRepo(db)
		trackedUsers, err := relationshipRepo.GetTrackedUsers(r.Context(), user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get tracked users")
			return
		}

		utils.RespondJSON(w, http.StatusOK, trackedUsers)
	}
}
