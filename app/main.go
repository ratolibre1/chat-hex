package main

import (
	"chat-hex/config"
	"chat-hex/modules/mongodb"
	"context"
	"encoding/json"
	"fmt"

	"chat-hex/api"
	authController "chat-hex/api/v1/auth"
	chatroomsController "chat-hex/api/v1/chatrooms"
	messagesController "chat-hex/api/v1/messages"
	preloadController "chat-hex/api/v1/preload"
	"chat-hex/api/v1/rabbit/requests"
	usersController "chat-hex/api/v1/users"
	authService "chat-hex/business/auth"
	chatroomsService "chat-hex/business/chatrooms"
	commandsService "chat-hex/business/commands"
	emitterService "chat-hex/business/emitter"
	listenerService "chat-hex/business/listener"
	messagesService "chat-hex/business/messages"
	preloadService "chat-hex/business/preload"
	usersService "chat-hex/business/users"
	chatroomsRepository "chat-hex/modules/chatrooms"
	messagesRepository "chat-hex/modules/messages"
	preloadRepository "chat-hex/modules/preload"
	usersRepository "chat-hex/modules/users"

	"os"
	"os/signal"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	//initiate emitter
	emitterService := emitterService.NewService()

	//initiate preload
	preloadRepo := preloadRepository.NewMongoDBRepository(dbConnection)
	preloadService := preloadService.NewService(preloadRepo)
	preloadController := preloadController.NewController(preloadService)

	//initiate chatrooms
	chatroomsRepo := chatroomsRepository.NewMongoDBRepository(dbConnection)
	chatroomsService := chatroomsService.NewService(chatroomsRepo)
	chatroomsController := chatroomsController.NewController(chatroomsService)

	//initiate users
	usersRepo := usersRepository.NewMongoDBRepository(dbConnection)
	usersService := usersService.NewService(usersRepo, chatroomsService)
	usersController := usersController.NewController(usersService)

	//initiate auth
	authService := authService.NewService()
	authController := authController.NewController(authService, usersService)


	//initiate messages
	messagesRepo := messagesRepository.NewMongoDBRepository(dbConnection)
	messagesService := messagesService.NewService(messagesRepo)

	//configure RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
    "chatQueue",
    false,
    false,
    false,
    false,
    nil,
	)
	if err != nil {
			panic(err)
	}
	
	//initiate commands
	commandsService := commandsService.NewService(messagesService, emitterService, ch)
	messagesController := messagesController.NewController(messagesService, commandsService)

	//initiate listener
	listenerService := listenerService.NewService()

	//register paths
	api.RegisterPaths(e, authController, preloadController, usersController, chatroomsController, messagesController)

	//setup the RabbitMQ consumer
	msgs, err := listenerService.ConsumeMessages(ch, "chatQueue")
	if err != nil {
		panic(err)
	}

	//process queue
	go func() {
		for msg := range msgs {

			var stockResponse requests.StockRequestResponse

			err := json.Unmarshal(msg.Body, &stockResponse)
			if err != nil {
					log.Printf("Error al deserializar JSON: %v", err)
					continue
			}

			err = commandsService.AsyncStockCommand(stockResponse.StockCode, stockResponse.Chatroom)
			if err != nil{
				log.Info("error processing stock command")
			}
		}
	}()

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
