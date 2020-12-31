package usecase

import (
	"context"
	"time"

	"github.com/blacknvcone/opdrewski/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type articleUseCase struct {
	articleRepo    domain.ArticleRepository
	contextTimeout time.Duration
}

func NewArticleUseCase(a domain.ArticleRepository, timeout time.Duration) domain.ArticleUseCase {
	return &articleUseCase{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *articleUseCase) Fetch(ctx context.Context, filter bson.M) (res []*domain.Article, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *articleUseCase) Store(ctx context.Context, ar *domain.Article) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	ar.ID = primitive.NewObjectID()
	ar.CreatedAt = time.Now()
	ar.UpdatedAt = time.Now()

	res, err := a.articleRepo.Store(ctx, ar)
	if err != nil {
		return nil, err
	}

	return res, nil

}
