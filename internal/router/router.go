package router

import (
	"gst-api/internal/handler"
	"gst-api/internal/middleware"
	"gst-api/internal/repository"

	"github.com/gin-gonic/gin"
)

func Setup(
	syncRepo *repository.SyncRepository,
	userRepo *repository.UserRepository,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	syncHandler := handler.NewSyncHandler(syncRepo)
	userHandler := handler.NewUserHandler(userRepo, jwtSecret)

	// ── Public routes ─────────────────────────────────────────────────────
	r.POST("/generate-token", userHandler.GenerateToken)
	r.POST("/api/v1/users", userHandler.CreateUser)

	// ── JWT protected: master data push & update ──────────────────────────
	masterData := r.Group("/api/v1/master-data")
	masterData.Use(middleware.JWTAuth(jwtSecret))
	{
		masterData.POST("/push", syncHandler.CreateSync)
		masterData.PATCH("/update", syncHandler.UpdateSync)
	}

	return r
}
