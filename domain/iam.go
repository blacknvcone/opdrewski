package domain

import (
	"context"
	"time"

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

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type IAMUseCase interface {
	//Authentication(ctx context.Context, req bson.M)
	//RefreshToken(ctx context.Context, req bson.M)
	AddUser(ctx context.Context, user *IAMUser) (interface{}, error)
	Authentication(ctx context.Context, email string, password string) (interface{}, error)
}

type IAMRepository interface {
	Fetch(ctx context.Context, filter bson.M) (*IAMUser, error)
	StoreUser(ctx context.Context, user *IAMUser) (interface{}, error)
}
