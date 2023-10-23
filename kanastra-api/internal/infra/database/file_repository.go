package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kanastra-api/internal/entity"
	"time"
)

type FileRepository struct {
	client *mongo.Client
}

func NewFileRepository(client *mongo.Client) *FileRepository {
	return &FileRepository{
		client: client,
	}
}

func (m *FileRepository) Save(file *entity.File) error {
	coll := m.client.Database("kanastra").Collection("files")

	result, err := coll.InsertOne(
		context.Background(),
		file,
	)

	file.ID = result.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		return err
	}

	return nil
}

func (m *FileRepository) FindAll() ([]entity.File, error) {
	coll := m.client.Database("kanastra").Collection("files")

	cursor, err := coll.Find(
		context.Background(),
		primitive.M{},
	)

	if err != nil {
		return nil, err
	}

	var files []entity.File

	err = cursor.All(context.Background(), &files)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (m *FileRepository) Update(id string) (*entity.File, error) {
	coll := m.client.Database("kanastra").Collection("files")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	after := options.After

	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	result := coll.FindOneAndUpdate(
		context.Background(),
		primitive.M{"_id": objectID},
		primitive.M{"$set": primitive.M{"processed": true, "processed_at": primitive.NewDateTimeFromTime(time.Now())}},
		opts,
	)

	var file entity.File

	err = result.Decode(&file)

	if err != nil {
		return nil, err
	}

	return &file, nil
}
