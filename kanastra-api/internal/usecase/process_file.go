package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"kanastra-api/internal/entity"
	"kanastra-api/pkg/events"
	"log"
	"os"
	"strconv"
	"sync"
)

type ProcessFileInputDTO struct {
	ID   string
	Path string
}

type FileUpdatedEvent struct {
	Message string      `json:"message"`
	Data    entity.File `json:"data"`
}

type ProcessFileUseCase struct {
	FileRepository   entity.FileRepositoryInterface
	BilletRepository entity.BilletRepositoryInterface
	FileEvent        events.EventInterface
	EventDispatcher  events.EventDispatcherInterface
}

func NewProcessFileUseCase(
	fileRepository entity.FileRepositoryInterface,
	billetRepository entity.BilletRepositoryInterface,
	fileEvent events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *ProcessFileUseCase {
	return &ProcessFileUseCase{
		FileRepository:   fileRepository,
		BilletRepository: billetRepository,
		FileEvent:        fileEvent,
		EventDispatcher:  eventDispatcher,
	}
}

func (u *ProcessFileUseCase) Worker(ctx context.Context, dst chan []string, src chan []string) {
	for {
		select {
		case url, ok := <-src: // Checar se o canal está fechado
			if !ok {
				return
			}

			go u.CreateBillet(url)

			dst <- url
		case <-ctx.Done(): // Checar se o contexto foi cancelado
			return
		}
	}
}

func (u *ProcessFileUseCase) SendEmail(name, email string) {
	fmt.Println("Sending email to", name, "at", email)
}

func (u *ProcessFileUseCase) CreateBillet(url []string) {
	fmt.Println("Creating billet for", url[0])

	u.SendEmail(url[0], url[2])
}

func (u *ProcessFileUseCase) SaveBillet(url []string) {
	debtAmount, err := strconv.ParseFloat(url[3], 64)

	if err != nil {
		fmt.Println("Error converting string to float64")
		fmt.Println(url[3])
		fmt.Println(err)
		return
	}

	billet := entity.NewBillet(url[0], url[1], url[2], url[4], url[5], debtAmount)

	err = u.BilletRepository.Save(billet)

	if err != nil {
		log.Fatal(err)
	}

	u.SendEmail(url[0], url[2])
}

func (u *ProcessFileUseCase) SaveBillets(urls [][]string) {
	for _, url := range urls {
		u.SaveBillet(url)
	}
}

func (u *ProcessFileUseCase) Execute(input *ProcessFileInputDTO) error {
	ctx := context.Background()
	file, err := os.Open(input.Path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Cria os canais
	src := make(chan []string)
	out := make(chan []string)

	// Cria o WaitGroup
	var wg sync.WaitGroup

	// Declara os workers
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.Worker(ctx, out, src)
		}()
	}

	// Lê o arquivo e envia para o canal src
	go func() {
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			src <- record // Envia o registro para o canal src para os workers consumirem
		}
		close(src) // Fecha o canal src
	}()

	// Espera os workers terminarem
	go func() {
		wg.Wait() // Espera os workers terminarem

		fileUpdated, err := u.FileRepository.Update(input.ID)

		if err != nil {
			return
		}

		u.FileEvent.SetPayload(FileUpdatedEvent{
			Message: "process-finished",
			Data:    *fileUpdated,
		})

		err = u.EventDispatcher.Dispatch(u.FileEvent)
		close(out) // Fecha o canal out
	}()

	billets := make([][]string, 0)
	for url := range out {
		billets = append(billets, url)
	}

	go u.SaveBillets(billets)

	if err != nil {
		return err
	}

	return nil
}
