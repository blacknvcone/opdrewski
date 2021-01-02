package http

import (
	"net/http"

	"github.com/blacknvcone/opdrewski/common/logger"
	"github.com/blacknvcone/opdrewski/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	AUseCase domain.ArticleUseCase
	log      logger.LogInfoFormat
}

func NewArticleHandler(router *gin.Engine, aUse domain.ArticleUseCase, logger logger.LogInfoFormat) {
	handler := &ArticleHandler{
		AUseCase: aUse,
		log:      logger,
	}

	router.GET("/articles", handler.FetchArticle)
	router.POST("/article", handler.StoreArticle)
}

func (a *ArticleHandler) FetchArticle(g *gin.Context) {
	ctx := g.Request.Context()
	listAr, err := a.AUseCase.Fetch(ctx, bson.M{})
	if err != nil {
		a.log.Info(err.Error())
		g.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "OK",
		"data":    listAr,
	})
}

func (a *ArticleHandler) StoreArticle(g *gin.Context) {
	article := domain.Article{}
	err := g.Bind(&article)
	if err != nil {
		a.log.Info(err.Error())
		g.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	ctx := g.Request.Context()
	res, err := a.AUseCase.Store(ctx, &article)
	if err != nil {
		a.log.Info(err.Error())
		g.JSON(http.StatusBadGateway, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	g.JSON(http.StatusCreated, res)
}
