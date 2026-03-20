package models

import "time"

type Geofence struct {
	ID            string
	Name          string
	CaregiverID   string
	TrackedUserID string
	Radius        int
	Latitude      float64
	Longitude     float64
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateGeofenceRequest struct {
	Name          string  `json:"name" validate:"required"`
	TrackedUserID string  `json:"tracked_user_id" validate:"required"`
	CaregiverID   string  `json:"caregiver_id" validate:"required"`
	Radius        int     `json:"radius" validate:"required,min=10,max=10000"`
	Latitude      float64 `json:"latitude" validate:"required"`
	Longitude     float64 `json:"longitude" validate:"required"`
}

type ListGeofencesForUserRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required,oneof=caregiver tracked"`
	Filter string `json:"filter" validate:"required,oneof=active inactive all"`
}

type UpdateGeofenceRequest struct {
	GeofenceID string  `json:"geofence_id" validate:"required"`
	Name       string  `json:"name" validate:"required"`
	Radius     int     `json:"radius" validate:"required,min=10,max=10000"`
	Latitude   float64 `json:"latitude" validate:"required"`
	Longitude  float64 `json:"longitude" validate:"required"`
	IsActive   bool    `json:"is_active" validate:"required"`
}

type DeleteGeofenceRequest struct {
	GeofenceID string `json:"geofence_id" validate:"required"`
}
