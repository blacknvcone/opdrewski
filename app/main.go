package main

import (
	"context"
	"log"
	"time"

	_articleDelivery "github.com/blacknvcone/opdrewski/article/delivery/http"
	_articleRepo "github.com/blacknvcone/opdrewski/article/repository/mongo"
	_articleUseCase "github.com/blacknvcone/opdrewski/article/usecase"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	cOpt := options.Client().ApplyURI("mongodb://localhost:27017")
	cli, err := mongo.NewClient(cOpt)
	if err != nil {
		log.Fatal(err)
	}

	err = cli.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//Ping
	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect into mongo instance :", err)
	}

	log.Println("Connected into databases !")

	e := echo.New()
	articleRepo := _articleRepo.NewMgoArticleRepository(cli)
	contextTimeout := time.Duration(10 * time.Second)
	articleUsecase := _articleUseCase.NewArticleUseCase(articleRepo, contextTimeout)
	_articleDelivery.NewArticleHandler(e, articleUsecase)

	log.Println(e.Start(":9090"))
}
