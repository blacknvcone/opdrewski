package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blacknvcone/opdrewski/domain"
	helper "github.com/blacknvcone/opdrewski/iam/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type iamUseCase struct {
	IAMRepository  domain.IAMRepository
	contextTimeout time.Duration
}

type claims struct {
	UUID string
	jwt.StandardClaims
}

func NewIAMUseCase(repository domain.IAMRepository, timeout time.Duration) domain.IAMUseCase {
	return &iamUseCase{
		IAMRepository:  repository,
		contextTimeout: timeout,
	}
}

func (u *iamUseCase) ValidateTokenHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Retrieving Token
		ts := c.Request.Header.Get("Authorization")
		if ts == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("lontongbalap"), nil
		})

		if token != nil && err == nil {
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Invalid Token",
			})
			c.Abort()
			return
		}

	}
}

func (u *iamUseCase) GenerateToken(ctx context.Context, uid string) (*domain.IAMToken, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	uuidG := uuid.New().String()
	expTime := time.Now().Add(5 * time.Minute)
	claims := &claims{}
	claims.UUID = uuidG
	claims.StandardClaims.ExpiresAt = expTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("lontongbalap")) //TODO : Make Secret to be params set
	if err != nil {
		return nil, err
	}

	iamt := &domain.IAMToken{}
	iamt.ID = primitive.NewObjectID()
	iamt.Expires = expTime.Unix()
	iamt.AccessToken = tokenString
	iamt.UUID = uuidG
	iamt.UserID = uid

	_, err = u.IAMRepository.StoreToken(ctx, iamt)
	if err != nil {
		return nil, err
	}

	return iamt, nil
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

	//Generate Token
	token, err := u.GenerateToken(ctx, res.ID.Hex())
	if err != nil {
		return nil, err
	}

	return bson.M{"AccessToken": token.AccessToken, "Expires": token.Expires}, nil

}
