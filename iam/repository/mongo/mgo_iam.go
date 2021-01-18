package mongo

import (
	"context"
	"os"

	"github.com/blacknvcone/opdrewski/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collIAMUser = "IAM_User"
const collIAMToken = "IAM_Token"

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

func (m *mgoIAMRepository) StoreToken(ctx context.Context, iamt *domain.IAMToken) (interface{}, error) {
	coll := m.db.Database(os.Getenv("MONGO_DB")).Collection(collIAMToken)
	filter := bson.M{
		"uid": iamt.UID,
	}

	updated := bson.M{
		"$set": bson.M{
			"accesstoken": iamt.AccessToken,
			"expires":     iamt.Expires,
		},
	}

	res := coll.FindOneAndUpdate(ctx, filter, updated)
	if res.Err() != nil {
		coll.InsertOne(ctx, iamt)
		return iamt, nil
	} else {
		return res.Decode(iamt), nil
	}

}
