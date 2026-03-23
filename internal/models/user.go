package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

// User represents a system user who can generate tokens for the master-data API.
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"   json:"id,omitempty"`
	Name      string             `bson:"name"            json:"name"            binding:"required"`
	Username  string             `bson:"username"        json:"username"        binding:"required"`
	Password  string             `bson:"password"        json:"-"` // never returned in JSON
	Email     string             `bson:"email"           json:"email"           binding:"required"`
	Contact   string             `bson:"contact"         json:"contact"         binding:"required"`
	IsActive  bool               `bson:"isActive"        json:"isActive"`
	CreatedAt time.Time          `bson:"createdAt"       json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"       json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Request / Response DTOs
// ---------------------------------------------------------------------------

// CreateUserRequest is the payload for POST /api/v1/users
type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"    binding:"required"`
	Contact  string `json:"contact"  binding:"required"`
}

// GenerateTokenRequest is the payload for POST /generate-token
type GenerateTokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse is returned after successful token generation
type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      UserInfo  `json:"user"`
}

// UserInfo is a safe public subset of User
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
}
