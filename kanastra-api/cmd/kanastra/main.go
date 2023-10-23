package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kanastra-api/configs"
	"kanastra-api/internal/event/handler"
	"kanastra-api/internal/infra/web/webserver"
	"kanastra-api/pkg/events"
)

func main() {
	env, err := configs.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?retryWrites=true&w=majority", env.MongoUser, env.MongoPassword, env.MongoHost)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		panic(err)
	}

	sseChannel := new(chan string)

	eventDispatcher := events.NewEventDispatcher()
	err = eventDispatcher.Register("FileEvent", &handler.FileEventHandler{
		SseChannel: sseChannel,
	})

	if err != nil {
		return
	}

	server := webserver.NewWebServer(env.ServerPort)
	webFileHandler := NewFileHandler(client, eventDispatcher)
	webEventsHandler := NewEventsHandler(sseChannel)
	server.AddHandler("/file", webFileHandler.Handle)
	server.AddHandler("/events", webEventsHandler.SendEvent)
	fmt.Println("Starting web server on port", env.ServerPort)
	server.Start()
}
