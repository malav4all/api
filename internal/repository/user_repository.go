package repository

import (
	"context"
	"time"

	"gst-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ColUser = "users"

// UserRepository handles all DB operations for the User module.
type UserRepository struct {
	db *mongo.Database
}

// NewUserRepository creates a UserRepository and ensures indexes.
func NewUserRepository(db *mongo.Database) (*UserRepository, error) {
	r := &UserRepository{db: db}
	// if err := r.ensureUserIndexes(context.Background()); err != nil {
	// 	return nil, err
	// }
	return r, nil
}

func (r *UserRepository) ensureUserIndexes(ctx context.Context) error {
	// Unique index on username and email
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uq_username"),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uq_email"),
		},
	}
	for _, idx := range indexes {
		_, err := r.db.Collection(ColUser).Indexes().CreateOne(ctx, idx)
		if err != nil && !mongo.IsDuplicateKeyError(err) {
			// Non-fatal — log handled by caller
			_ = err
		}
	}
	return nil
}

// CreateUser inserts a new user. Password must already be hashed before calling.
func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	_, err := r.db.Collection(ColUser).InsertOne(ctx, user)
	return err
}

// FindByUsername returns the user matching the given username (active only).
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.Collection(ColUser).FindOne(ctx,
		bson.M{"username": username, "isActive": true},
	).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &user, err
}
