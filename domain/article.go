package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Title     string             `json:"title,omitempty" validate:"required" bson:"title"`
	Content   string             `json:"content,omitempty" validate:"required" bson:"content"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
}

type ArticleUseCase interface {
	Fetch(ctx context.Context, filter bson.M) ([]*Article, error)
	Store(ctx context.Context, ar *Article) (interface{}, error)
	// GetByID(id string) (Article, error)

}

type ArticleRepository interface {
	Fetch(ctx context.Context, filter bson.M) ([]*Article, error)
	Store(ctx context.Context, ar *Article) (interface{}, error)
	// GetByID(id string) (Article, error)

}
