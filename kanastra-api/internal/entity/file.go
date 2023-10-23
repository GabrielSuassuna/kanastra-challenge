package entity

import (
	"errors"
	"strings"
	"time"
)

type File struct {
	ID          string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string    `bson:"name" json:"name"`
	Processed   bool      `bson:"processed" json:"processed"`
	ProcessedAt time.Time `bson:"processed_at" json:"processed_at,omitempty"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

func NewFile(name string) *File {
	createdAt := time.Now()
	updatedAt := time.Now()

	return &File{
		Name:      name,
		Processed: false,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (f *File) IsValid() error {
	if f.Name == "" {
		return errors.New("name is required")
	}

	if strings.Split(f.Name, ".")[1] != "csv" {
		return errors.New("mime type must be csv")
	}
	return nil
}
