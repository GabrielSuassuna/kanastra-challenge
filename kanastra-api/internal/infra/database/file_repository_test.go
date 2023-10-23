package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"kanastra-api/internal/entity"
	"reflect"
	"testing"
)

func TestFileRepository_Save(t *testing.T) {
	type fields struct {
		client *mongo.Client
	}
	type args struct {
		file *entity.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"should save a file",
			fields{client: nil},
			args{file: &entity.File{Name: "file.csv"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &FileRepository{
				client: tt.fields.client,
			}
			if err := m.Save(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewFileRepository(t *testing.T) {
	type args struct {
		client *mongo.Client
	}
	tests := []struct {
		name string
		args args
		want *FileRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileRepository(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
