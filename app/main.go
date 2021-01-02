package main

import (
	"context"
	"log"
	"os"
	"time"

	_articleDelivery "github.com/blacknvcone/opdrewski/article/delivery/http"
	_articleRepo "github.com/blacknvcone/opdrewski/article/repository/mongo"
	_articleUseCase "github.com/blacknvcone/opdrewski/article/usecase"
	_logger "github.com/blacknvcone/opdrewski/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	//Loading Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Initiate Logger
	logger, _ := _logger.NewLogger()

	cOpt := options.Client().ApplyURI("mongodb://" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT"))
	cli, err := mongo.NewClient(cOpt)
	if err != nil {
		logger.Fatal(err)
	}

	err = cli.Connect(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	//Ping
	err = cli.Ping(context.Background(), readpref.Primary())
	if err != nil {
		logger.Fatal("Couldn't connect into mongo instance :", err)
	}

	logger.Info("Connected into databases !")

	router := gin.New()
	contextTimeout := time.Duration(10 * time.Second)
	articleRepo := _articleRepo.NewMgoArticleRepository(cli)
	articleUsecase := _articleUseCase.NewArticleUseCase(articleRepo, contextTimeout)
	_articleDelivery.NewArticleHandler(router, articleUsecase, logger)

	router.Run(":" + os.Getenv("SERVER_PORT"))

}
