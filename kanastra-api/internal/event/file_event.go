package event

import "time"

type FileEvent struct {
	Name    string
	Payload interface{}
}

func NewFileEvent() *FileEvent {
	return &FileEvent{
		Name: "FileEvent",
	}
}

func (e *FileEvent) GetName() string {
	return e.Name
}

func (e *FileEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *FileEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *FileEvent) GetDateTime() time.Time {
	return time.Now()
}
