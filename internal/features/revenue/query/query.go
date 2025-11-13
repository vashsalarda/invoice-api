package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/revenue/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultRevenueQuery struct{}

func (c *DefaultRevenueQuery) CollectionName() string {
	return "revenues"
}

//go:generate mockgen -destination=../mocks/query/mock_invoice_query.go -package=query invoice-api/internal/features/invoice/query RevenueQuery
type RevenueQuery interface {
	GetItemsByQuery() ([]model.RevenueDTO, error)
	GetItemByID(id string) (*model.RevenueDTO, error)
}

func (c *DefaultRevenueQuery) GetItemsByQuery() ([]model.RevenueDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	items := make([]model.RevenueDTO, 0, 100) 
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "createdAt", Value: 1}})
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return items, err
	}

	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &items); err != nil {
		return items, err
	}

	return items, nil
}

func (c *DefaultRevenueQuery) GetItemByID(id string) (*model.RevenueDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item model.RevenueDTO
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &item, nil
}