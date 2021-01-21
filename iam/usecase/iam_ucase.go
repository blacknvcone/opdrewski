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

		//Verify
		token, err := verifyToken(ts)

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

func verifyToken(ts string) (*jwt.Token, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		//does this token conform to "SigningMethodHMAC" ?
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("lontongbalap"), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (u *iamUseCase) ExtractSession(ctx context.Context, ts string) (*domain.IAMUser, error) {
	//Verify
	token, err := verifyToken(ts)
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims

	if ok && token.Valid {
		q := bson.M{
			"uuid": claims["UUID"],
		}

		iamtoken, err := u.IAMRepository.FetchSession(ctx, q)
		if err != nil {
			return nil, err
		}

		uid, err := primitive.ObjectIDFromHex(iamtoken.UserID)
		if err != nil {
			return nil, err
		}

		q = bson.M{
			"_id": uid,
		}

		iamuser, err := u.IAMRepository.Fetch(ctx, q)
		if err != nil {
			return nil, err
		}

		return iamuser, nil

	}
	return nil, err
}

func (u *iamUseCase) GenerateToken(ctx context.Context, uuid string, expired int64) (string, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	expTime := time.Now().Add(5 * time.Minute)
	claims := &claims{}
	claims.UUID = uuid
	claims.StandardClaims.ExpiresAt = expTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("lontongbalap")) //TODO : Make Secret to be params set
	if err != nil {
		return "", 0, err
	}

	return tokenString, expired, nil
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

	iamuser, err := u.IAMRepository.Fetch(ctx, filter)
	if err != nil {
		return nil, err
	}

	//Generate Session
	uuidGen := uuid.New().String()
	expSess := time.Now().Add(5 * time.Minute)

	iams := &domain.IAMSession{}
	iams.ID = primitive.NewObjectID()
	iams.Expires = expSess.Unix()
	iams.UUID = uuidGen
	iams.UserID = iamuser.ID.Hex()

	_, err = u.IAMRepository.StoreSession(ctx, iams)
	if err != nil {
		return "", err
	}

	//Generate Token
	expToken := time.Now().Add(5 * time.Minute)
	token, exp, err := u.GenerateToken(ctx, uuidGen, expToken.Unix())
	if err != nil {
		return nil, err
	}

	return bson.M{"AccessToken": token, "Expires": exp}, nil

}
