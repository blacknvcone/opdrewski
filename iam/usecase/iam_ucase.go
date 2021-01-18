package usecase

import (
	"context"
	"time"

	"github.com/blacknvcone/opdrewski/domain"
	helper "github.com/blacknvcone/opdrewski/iam/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type iamUseCase struct {
	IAMRepository  domain.IAMRepository
	contextTimeout time.Duration
}

func NewIAMUseCase(repository domain.IAMRepository, timeout time.Duration) domain.IAMUseCase {
	return &iamUseCase{
		IAMRepository:  repository,
		contextTimeout: timeout,
	}
}

func (u *iamUseCase) AddUser(ctx context.Context, user *domain.IAMUser) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.Password = helper.GetMD5Hash(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	res, err := u.IAMRepository.StoreUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *iamUseCase) Authentication(ctx context.Context, email string, password string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	filter := bson.M{
		"email":    email,
		"password": helper.GetMD5Hash(password),
	}

	res, err := u.IAMRepository.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	//TODO :
	//Add JWT Token Here Soon
	return res, nil

}
