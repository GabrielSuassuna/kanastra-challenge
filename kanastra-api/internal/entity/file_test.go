package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestFile_IsValid(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		MimeType  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"should return error when name is empty", fields{Name: "", MimeType: "csv"}, true},
		{"should return error when mime type is empty", fields{Name: "file.csv", MimeType: ""}, true},
		{"should return error when mime type is not csv", fields{Name: "file.csv", MimeType: "pdf"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if err := f.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewFile(t *testing.T) {
	type args struct {
		name     string
		mimeType string
	}
	tests := []struct {
		name string
		args args
		want *File
	}{
		{"should return a new file", args{name: "file.csv"}, &File{Name: "file.csv"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFile(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
