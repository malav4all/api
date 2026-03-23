package handler

import (
	"net/http"

	"gst-api/internal/models"
	"gst-api/internal/repository"
	"gst-api/pkg/jwt"
	"gst-api/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler handles user creation and token generation.
type UserHandler struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserHandler(repo *repository.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{repo: repo, jwtSecret: jwtSecret}
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/v1/users  — create a new user
// ─────────────────────────────────────────────────────────────────────────────

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Hash password before storing
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalServerError(c, "failed to hash password")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
		Contact:  req.Contact,
	}

	if err := h.repo.CreateUser(c.Request.Context(), user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "username or email already exists",
			})
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	// Return user without password
	response.Created(c, "user created successfully", models.UserInfo{
		ID:       user.ID.Hex(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Contact:  user.Contact,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /generate-token  — authenticate and get JWT
// ─────────────────────────────────────────────────────────────────────────────

func (h *UserHandler) GenerateToken(c *gin.Context) {
	var req models.GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Find user by username
	user, err := h.repo.FindByUsername(c.Request.Context(), req.Username)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid username or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid username or password",
		})
		return
	}

	// Generate JWT
	token, expiresAt, err := jwt.Generate(user.ID.Hex(), user.Username, h.jwtSecret)
	if err != nil {
		response.InternalServerError(c, "failed to generate token")
		return
	}

	response.Success(c, "token generated successfully", models.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: models.UserInfo{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Contact:  user.Contact,
		},
	})
}
