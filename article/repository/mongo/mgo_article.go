package mongo

import (
	"context"

	"github.com/blacknvcone/opdrewski/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type mgoArticleRepository struct {
	db      string
	session *mongo.Session
}

func NewMgoArticleRepository(ctx ctx) *domain.ArticleRepository {

}

func (m *mgoArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {

}
