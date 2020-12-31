package http

import (
	"net/http"

	"github.com/blacknvcone/opdrewski/domain"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ArticleHandler struct {
	AUseCase domain.ArticleUseCase
}

func NewArticleHandler(e *echo.Echo, aUse domain.ArticleUseCase) {
	handler := &ArticleHandler{
		AUseCase: aUse,
	}
	e.GET("/articles", handler.FetchArticle)
	e.POST("/article", handler.StoreArticle)
}

func (a *ArticleHandler) FetchArticle(c echo.Context) error {
	ctx := c.Request().Context()
	listAr, err := a.AUseCase.Fetch(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusBadGateway, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, listAr)
}

func (a *ArticleHandler) StoreArticle(c echo.Context) (err error) {
	article := domain.Article{}
	err = c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	ctx := c.Request().Context()
	res, err := a.AUseCase.Store(ctx, &article)
	if err != nil {
		return c.JSON(http.StatusBadGateway, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}
