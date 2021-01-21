package domain

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAMUser struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at"`
}

type IAMToken struct {
	ID          primitive.ObjectID `json:"_id,omitempty"  bson:"_id"`
	UUID        string             `json:",omitempty"`
	UserID      string             `json:",omitempty"`
	AccessToken string             `json:",omitempty"`
	Expires     int64              `json:",omitempty"`
}

type IAMUseCase interface {
	AddUser(ctx context.Context, user *IAMUser) (interface{}, error)
	Authentication(ctx context.Context, email string, password string) (interface{}, error)
	GenerateToken(ctx context.Context, uid string) (*IAMToken, error)
	ValidateTokenHTTP() gin.HandlerFunc
}

type IAMRepository interface {
	Fetch(ctx context.Context, filter bson.M) (*IAMUser, error)
	StoreUser(ctx context.Context, user *IAMUser) (interface{}, error)
	StoreToken(ctx context.Context, token *IAMToken) (interface{}, error)
}
