package domain

import (
	"context"
	"time"
)

type Article struct {
	ID      string `json:"_id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	//Author    Author    `json:"author"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleUseCase interface {
	Fetch(ctx context.Context, cursor string) ([]Article, string, error)
	GetByID(ctx context.Context, id string) (Article, error)
	Store(ctx context.Context, ar *Article) error
}

type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string) ([]Article, string, error)
	GetByID(ctx context.Context, id string) (Article, error)
	Store(ctx context.Context, ar *Article) error
}
