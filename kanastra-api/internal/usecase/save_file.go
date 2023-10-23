package usecase

import (
	"kanastra-api/internal/entity"
	"kanastra-api/pkg/events"
)

type SaveFileInputDTO struct {
	FileName string
	Path     string
}

type SaveFileOutputDTO struct {
	File entity.File
}

type SaveFileUseCase struct {
	FileRepository  entity.FileRepositoryInterface
	FileEvent       events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewSaveFileUseCase(
	fileRepository entity.FileRepositoryInterface,
	fileEvent events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *SaveFileUseCase {
	return &SaveFileUseCase{
		FileRepository:  fileRepository,
		FileEvent:       fileEvent,
		EventDispatcher: eventDispatcher,
	}
}

func (u *SaveFileUseCase) Execute(input *SaveFileInputDTO) (*SaveFileOutputDTO, error) {
	file := entity.NewFile(input.FileName)

	err := u.FileRepository.Save(file)

	if err != nil {
		return nil, err
	}

	output := &SaveFileOutputDTO{
		File: *file,
	}

	u.FileEvent.SetPayload("created-file")

	err = u.EventDispatcher.Dispatch(u.FileEvent)

	if err != nil {
		return nil, err
	}

	return output, nil
}
