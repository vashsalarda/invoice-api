package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/customer/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultQuery struct{}

func (c *DefaultQuery) CollectionName() string {
	return "customers"
}

type Query interface {
	GetCustomerByQuery() ([]model.CustomerDTO, error)
	GetCustomerByID(id string) (*model.CustomerDTO, error)
}


func (c *DefaultQuery) GetCustomerByQuery() ([]model.CustomerDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]model.CustomerDTO, 0, 100)
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *DefaultQuery) GetCustomerByID(id string) (*model.CustomerDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item model.CustomerDTO
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &item, nil
}
