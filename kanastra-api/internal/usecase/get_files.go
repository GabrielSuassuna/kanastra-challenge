package usecase

import "kanastra-api/internal/entity"

type GetFilesInputDTO struct {
}

type GetFilesOutputDTO struct {
	Files []entity.File `json:"files"`
}

type GetFilesUseCase struct {
	FileRepository entity.FileRepositoryInterface
}

func NewGetFilesUseCase(fileRepository entity.FileRepositoryInterface) *GetFilesUseCase {
	return &GetFilesUseCase{
		FileRepository: fileRepository,
	}
}

func (u *GetFilesUseCase) Execute(input *GetFilesInputDTO) (*GetFilesOutputDTO, error) {
	files, err := u.FileRepository.FindAll()

	if err != nil {
		return nil, err
	}

	output := &GetFilesOutputDTO{
		Files: files,
	}

	return output, nil
}
