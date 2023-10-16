package main

import (
	"chat-hex/config"
	"chat-hex/modules/mongodb"
	"context"
	"fmt"

	"chat-hex/api"
	messagesController "chat-hex/api/v1/messages"
	usersController "chat-hex/api/v1/users"
	commandsService "chat-hex/business/commands"
	messagesService "chat-hex/business/messages"
	usersService "chat-hex/business/users"
	messagesRepository "chat-hex/modules/messages"
	usersRepository "chat-hex/modules/users"

	"os"
	"os/signal"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func newDatabaseConnection(config *config.AppConfig) *mongo.Database {
	uri := "mongodb://"

	if config.AppEnvironment == "prod" {
		uri = "mongodb+srv://"
	}

	if config.DbUsername != "" {
		uri = fmt.Sprintf("%s%v:%v@", uri, config.DbUsername, config.DbPassword)
	}

	if config.AppEnvironment == "prod" {
		uri = fmt.Sprintf("%s%v/factura?retryWrites=true&w=majority",
			uri,
			config.DbAddress,
		)
	} else {
		uri = fmt.Sprintf("%s%v:%v/?connect=direct",
			uri,
			config.DbAddress,
			config.DbPort,
		)
	}

	db, err := mongodb.NewDatabaseConnection(uri, config.DbName)

	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	//load config if available or set to default
	config := config.GetConfig()

	//initialize database connection based on given config
	dbConnection := newDatabaseConnection(config)

	//create echo http
	e := echo.New()

	//initiate users
	usersRepo := usersRepository.NewMongoDBRepository(dbConnection)
	usersService := usersService.NewService(usersRepo)
	usersController := usersController.NewController(usersService)

	//initiate commands
	commandsService := commandsService.NewService()

	//initiate messages
	messagesRepo := messagesRepository.NewMongoDBRepository(dbConnection)
	messagesService := messagesService.NewService(messagesRepo)
	messagesController := messagesController.NewController(messagesService, commandsService)

	//register paths
	api.RegisterPaths(e, usersController, messagesController)

	// run server
	go func() {
		address := fmt.Sprintf("localhost:%d", config.AppPort)

		if err := e.Start(address); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
