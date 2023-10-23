package web

import (
	"encoding/json"
	"fmt"
	"io"
	"kanastra-api/internal/entity"
	"kanastra-api/internal/usecase"
	"kanastra-api/pkg/events"
	"net/http"
	"os"
	"strings"
	"time"
)

type FileHandler struct {
	EventDispatcher events.EventDispatcherInterface
	FileRepository  entity.FileRepositoryInterface
	FileEvent       events.EventInterface
}

func NewFileHandler(
	eventDispatcher events.EventDispatcherInterface,
	fileRepository entity.FileRepositoryInterface,
	fileEvent events.EventInterface,
) *FileHandler {
	return &FileHandler{
		EventDispatcher: eventDispatcher,
		FileRepository:  fileRepository,
		FileEvent:       fileEvent,
	}
}

func (h *FileHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Save(w, r)
	case http.MethodGet:
		h.Get(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *FileHandler) Save(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 << 20) // 128 MB is the maximum file size
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filenameSplit := strings.Split(header.Filename, ".")

	if len(filenameSplit) < 2 {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	if filenameSplit[1] != "csv" {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	filename := filenameSplit[0] + "_" + time.Now().Format("20060102150405") + "." + filenameSplit[1]

	f, err := os.OpenFile("./uploads/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the contents of the file to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saveFileUseCase := usecase.NewSaveFileUseCase(h.FileRepository, h.FileEvent, h.EventDispatcher)

	saveFileOutputDto, err := saveFileUseCase.Execute(
		&usecase.SaveFileInputDTO{
			FileName: filename,
			Path:     "./uploads/" + filename,
		})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	processFileUseCase := usecase.NewProcessFileUseCase(h.FileRepository, h.FileEvent, h.EventDispatcher)

	go func() {
		err = processFileUseCase.Execute(
			&usecase.ProcessFileInputDTO{
				ID:   saveFileOutputDto.File.ID,
				Path: "./uploads/" + filename,
			})

		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(saveFileOutputDto.File)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *FileHandler) Get(w http.ResponseWriter, r *http.Request) {
	getFilesUseCase := usecase.NewGetFilesUseCase(h.FileRepository)

	getFilesOutputDto, err := getFilesUseCase.Execute(&usecase.GetFilesInputDTO{})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(getFilesOutputDto.Files) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.NewEncoder(w).Encode(getFilesOutputDto.Files)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
