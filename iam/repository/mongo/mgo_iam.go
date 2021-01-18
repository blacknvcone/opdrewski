package mongo

import (
	"context"
	"os"

	"github.com/blacknvcone/opdrewski/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collIAMUser = "IAM_User"

type mgoIAMRepository struct {
	db *mongo.Client
}

func NewMgoIAMRepository(db *mongo.Client) domain.IAMRepository {
	return &mgoIAMRepository{db}
}

func (m *mgoIAMRepository) Fetch(ctx context.Context, filter bson.M) (IAMUser *domain.IAMUser, err error) {
	coll := m.db.Database(os.Getenv("MONGO_DB")).Collection(collIAMUser)
	res := coll.FindOne(ctx, filter)

	err = res.Decode(&IAMUser)
	if err != nil {
		return nil, bson.ErrDecodeToNil
	}

	return IAMUser, nil
}

func (m *mgoIAMRepository) StoreUser(ctx context.Context, iamu *domain.IAMUser) (interface{}, error) {
	coll := m.db.Database(os.Getenv("MONGO_DB")).Collection(collIAMUser)
	res, err := coll.InsertOne(ctx, iamu)
	if err != nil {
		return nil, mongo.MarshalError{}
	}

	return res, nil

}
