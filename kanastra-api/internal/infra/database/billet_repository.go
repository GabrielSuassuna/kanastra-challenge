package database

import (
	"context"
	"kanastra-api/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BilletRepository struct {
	client *mongo.Client
}

func NewBilletRepository(client *mongo.Client) *BilletRepository {
	return &BilletRepository{
		client: client,
	}
}

func (m *BilletRepository) Save(billet *entity.Billet) error {
	coll := m.client.Database("kanastra").Collection("billets")

	result, err := coll.InsertOne(
		context.Background(),
		billet,
	)

	billet.ID = result.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		return err
	}

	return nil
}
