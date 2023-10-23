package entity

type FileRepositoryInterface interface {
	Save(file *File) error
	FindAll() ([]File, error)
	Update(id string) (*File, error)
}
