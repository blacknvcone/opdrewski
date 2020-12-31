package mongo

import (
	"context"
	"log"
	"os"

	"github.com/blacknvcone/opdrewski/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collArticle = "articles"

type mgoArticleRepository struct {
	db *mongo.Client
}

func NewMgoArticleRepository(db *mongo.Client) domain.ArticleRepository {
	return &mgoArticleRepository{db}
}

func (m *mgoArticleRepository) Fetch(ctx context.Context, filter bson.M) (res []*domain.Article, err error) {
	coll := m.db.Database(os.Getenv("MONGO_DB")).Collection(collArticle)
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var article domain.Article
		err = cur.Decode(&article)
		if err != nil {
			log.Fatal("Error on Decoding the document : ", err)
		}
		res = append(res, &article)
	}

	return res, nil
}

func (m *mgoArticleRepository) Store(ctx context.Context, ar *domain.Article) (interface{}, error) {
	coll := m.db.Database(os.Getenv("MONGO_DB")).Collection(collArticle)
	res, err := coll.InsertOne(ctx, ar)
	if err != nil {
		log.Fatal("Error on Inserting the document : ", err)
	}

	return res, nil
}
