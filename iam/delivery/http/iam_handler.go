package http

import (
	"net/http"

	"github.com/blacknvcone/opdrewski/common/logger"
	"github.com/blacknvcone/opdrewski/domain"
	"github.com/gin-gonic/gin"
)

type ResponseError struct {
	Message string `json:"message"`
}

type IAMHandler struct {
	IAMUseCase domain.IAMUseCase
	log        logger.LogInfoFormat
}

func NewIAMHandler(router *gin.Engine, iamUse domain.IAMUseCase, logger logger.LogInfoFormat) {
	handler := &IAMHandler{
		IAMUseCase: iamUse,
		log:        logger,
	}

	router.POST("/iam/create", handler.CreateUser)
	router.POST("/iam/auth", handler.Auth)

}

func (i *IAMHandler) CreateUser(g *gin.Context) {
	user := domain.IAMUser{}
	err := g.Bind(&user)
	if err != nil {
		i.log.Info(err.Error())
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	ctx := g.Request.Context()
	res, err := i.IAMUseCase.AddUser(ctx, &user)

	if err != nil {
		i.log.Info(err.Error())
		g.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, res)

}

func (i *IAMHandler) Auth(g *gin.Context) {

	iam := domain.IAMUser{}
	err := g.Bind(&iam)
	if err != nil {
		i.log.Info(err.Error())
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	ctx := g.Request.Context()
	res, err := i.IAMUseCase.Authentication(ctx, iam.Email, iam.Password)
	if err != nil {
		i.log.Info(err.Error())
		g.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "OK",
		"data":    res,
	})

}
