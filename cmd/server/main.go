package main

import (
	"log"

	"gst-api/internal/config"
	"gst-api/internal/repository"
	"gst-api/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Set Gin mode
	gin.SetMode(cfg.GinMode)

	// 3. Connect to MongoDB
	log.Printf("[MAIN] Connecting to MongoDB at %s (db: %s)", cfg.MongoURI, cfg.MongoDBName)
	db, err := repository.Connect(cfg.MongoURI, cfg.MongoDBName)
	if err != nil {
		log.Fatalf("[MAIN] MongoDB connection failed: %v", err)
	}
	log.Println("[MAIN] MongoDB connected successfully")

	// 4. Initialise SyncRepository
	syncRepo, err := repository.New(db)
	if err != nil {
		log.Fatalf("[MAIN] SyncRepository init failed: %v", err)
	}

	// 5. Initialise UserRepository
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("[MAIN] UserRepository init failed: %v", err)
	}

	// 6. Setup HTTP router
	engine := router.Setup(syncRepo, userRepo, cfg.JWTSecret)

	// 7. Start server
	addr := ":" + cfg.ServerPort
	log.Printf("[MAIN] Server starting on %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("[MAIN] Server failed: %v", err)
	}
}
