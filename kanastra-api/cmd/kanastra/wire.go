//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"kanastra-api/internal/entity"
	"kanastra-api/internal/event"
	"kanastra-api/internal/infra/database"
	"kanastra-api/internal/infra/web"
	"kanastra-api/internal/usecase"
	"kanastra-api/pkg/events"
)

var setFileRepositoryDependency = wire.NewSet(
	database.NewFileRepository,
	wire.Bind(new(entity.FileRepositoryInterface), new(*database.FileRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewFileEvent,
	wire.Bind(new(events.EventInterface), new(*event.FileEvent)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setFileEvent = wire.NewSet(
	event.NewFileEvent,
	wire.Bind(new(events.EventInterface), new(*event.FileEvent)),
)

func NewCreateFileUseCase(client *mongo.Client, eventDispatcher events.EventDispatcherInterface) *usecase.SaveFileUseCase {
	wire.Build(
		setFileRepositoryDependency,
		setFileEvent,
		usecase.NewSaveFileUseCase,
	)
	return &usecase.SaveFileUseCase{}
}

func NewProcessFileUseCase(client *mongo.Client, eventDispatcher events.EventDispatcherInterface) *usecase.ProcessFileUseCase {
	wire.Build(
		setFileRepositoryDependency,
		setFileEvent,
		usecase.NewProcessFileUseCase,
	)
	return &usecase.ProcessFileUseCase{}
}

func NewGetFilesUseCase(client *mongo.Client) *usecase.GetFilesUseCase {
	wire.Build(
		setFileRepositoryDependency,
		usecase.NewGetFilesUseCase,
	)
	return &usecase.GetFilesUseCase{}
}

func NewFileHandler(client *mongo.Client, eventDispatcher events.EventDispatcherInterface) *web.FileHandler {
	wire.Build(
		setFileRepositoryDependency,
		setFileEvent,
		web.NewFileHandler,
	)
	return &web.FileHandler{}
}

func NewEventsHandler(sseChannel *chan string) *web.EventHandler {
	wire.Build(
		web.NewEventHandler,
	)

	return &web.EventHandler{}
}
