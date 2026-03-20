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

func CreateGeofence(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.CreateGeofenceRequest
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
			utils.RespondError(w, http.StatusForbidden, "Only caregivers can create geofences")
			return
		}

		// Check if caregiver owns tracked user
		relationshipRepo := database.NewRelationshipRepo(db)
		exists, err := relationshipRepo.CheckRelationshipExists(r.Context(), user.ID, req.TrackedUserID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to check relationship")
			return
		}
		if !exists {
			log.Println("Caregiver does not own tracked user")
			utils.RespondError(w, http.StatusForbidden, "You do not have permission to create a geofence for this user")
			return
		}

		geofenceRepo := database.NewGeofenceRepo(db)
		geofence, err := geofenceRepo.CreateGeofence(r.Context(), &req)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to create geofence")
			return
		}

		utils.RespondJSON(w, http.StatusCreated, geofence)
	}
}

func ListGeofencesForUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.ListGeofencesForUserRequest
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

		// Check if user is caregiver or tracked user
		if user.ID != req.UserID && user.Role != req.Role {
			log.Println("User is not allowed to list geofences for this user")
			utils.RespondError(w, http.StatusForbidden, "You do not have permission to list geofences for this user")
			return
		}

		geofenceRepo := database.NewGeofenceRepo(db)
		geofences, err := geofenceRepo.ListGeofencesForUser(r.Context(), &req)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to list geofences")
			return
		}

		utils.RespondJSON(w, http.StatusOK, geofences)
	}
}

func UpdateGeofence(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.UpdateGeofenceRequest
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
			utils.RespondError(w, http.StatusForbidden, "Only caregivers can update geofences")
			return
		}

		// Check if caregiver owns geofence
		geofenceRepo := database.NewGeofenceRepo(db)
		exists, err := geofenceRepo.CheckGeofenceOwnership(r.Context(), req.GeofenceID, user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to check geofence ownership")
			return
		}
		if !exists {
			log.Println("Caregiver does not own geofence")
			utils.RespondError(w, http.StatusForbidden, "You do not have permission to update this geofence")
			return
		}

		geofence, err := geofenceRepo.GetGeofenceByID(r.Context(), req.GeofenceID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to get geofence")
			return
		}

		// Update geofence fields
		geofence.Name = req.Name
		geofence.Radius = req.Radius
		geofence.Latitude = req.Latitude
		geofence.Longitude = req.Longitude
		geofence.IsActive = req.IsActive

		if err := geofenceRepo.UpdateGeofence(r.Context(), geofence); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to update geofence")
			return
		}

		log.Println("Geofence updated")
		utils.RespondJSON(w, http.StatusOK, geofence)
	}
}

func DeleteGeofence(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validate := validator.New()
		var req models.DeleteGeofenceRequest
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
			utils.RespondError(w, http.StatusForbidden, "Only caregivers can delete geofences")
			return
		}

		// Check if caregiver owns geofence
		geofenceRepo := database.NewGeofenceRepo(db)
		exists, err := geofenceRepo.CheckGeofenceOwnership(r.Context(), req.GeofenceID, user.ID)
		if err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to check geofence ownership")
			return
		}
		if !exists {
			log.Println("Caregiver does not own geofence")
			utils.RespondError(w, http.StatusForbidden, "You do not have permission to delete this geofence")
			return
		}

		if err := geofenceRepo.DeleteGeofence(r.Context(), req.GeofenceID); err != nil {
			log.Println(err)
			utils.RespondError(w, http.StatusInternalServerError, "Failed to delete geofence")
			return
		}

		log.Println("Geofence deleted")
		utils.RespondJSON(w, http.StatusOK, nil)
	}
}
