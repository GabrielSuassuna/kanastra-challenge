package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"kanastra-api/internal/entity"
	"kanastra-api/pkg/events"
	"log"
	"os"
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
	FileRepository  entity.FileRepositoryInterface
	FileEvent       events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewProcessFileUseCase(
	fileRepository entity.FileRepositoryInterface,
	fileEvent events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *ProcessFileUseCase {
	return &ProcessFileUseCase{
		FileRepository:  fileRepository,
		FileEvent:       fileEvent,
		EventDispatcher: eventDispatcher,
	}
}

func (u *ProcessFileUseCase) Worker(ctx context.Context, dst chan string, src chan []string) {
	for {
		select {
		case url, ok := <-src: // Checar se o canal está fechado
			if !ok {
				return
			}

			go u.SendEmail(url[0], url[1])

			dst <- url[0]
		case <-ctx.Done(): // Checar se o contexto foi cancelado
			return
		}
	}
}

func (u *ProcessFileUseCase) SendEmail(name, email string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Sua fatura está pronta!")
	m.SetBody("text/html", fmt.Sprintf("<h1>Olá %s, sua fatura está pronta!</h1>", name))

	d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")

	if err := d.DialAndSend(m); err != nil {

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
	out := make(chan string)

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

			src <- record // Envia o registro para o canal src
		}
		close(src) // Fecha o canal src
	}()

	// Espera os workers terminarem
	go func() {
		wg.Wait()  // Espera os workers terminarem
		close(out) // Fecha o canal out
	}()

	// Drena o canal out
	for _ = range out {
	}

	fileUpdated, err := u.FileRepository.Update(input.ID)

	if err != nil {
		return err
	}

	u.FileEvent.SetPayload(FileUpdatedEvent{
		Message: "process-finished",
		Data:    *fileUpdated,
	})

	err = u.EventDispatcher.Dispatch(u.FileEvent)

	if err != nil {
		return err
	}

	return nil
}
