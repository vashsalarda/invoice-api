package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/invoice/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultInvoiceQuery struct{}

func (c *DefaultInvoiceQuery) CollectionName() string {
	return "invoices"
}

//go:generate mockgen -destination=../mocks/query/mock_invoice_query.go -package=query invoice-api/internal/features/invoice/query InvoiceQuery
type InvoiceQuery interface {
	GetItemsByQuery() ([]model.InvoiceDTO, error)
	GetItemByID(id string) (*model.InvoiceDTO, error)
	GetLatestInvoices() ([]model.LatestInvoice, error)
}

func (c *DefaultInvoiceQuery) GetItemsByQuery() ([]model.InvoiceDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	items := make([]model.InvoiceDTO, 0, 100)
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
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

func (c *DefaultInvoiceQuery) GetItemByID(id string) (*model.InvoiceDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var item model.InvoiceDTO
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &item, nil
}

func (c *DefaultInvoiceQuery) GetLatestInvoices() ([]model.LatestInvoice, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$sort": bson.M{"date": -1},
		},
		{
			"$limit": 5,
		},
		{
			"$lookup": bson.M{
				"from":         "customers",
				"localField":   "customerId",
				"foreignField": "_id",
				"as":           "customer",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$customer",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$project": bson.M{
				"_id":      "$_id",
				"name":     "$customer.name",
				"imageUrl": "$customer.imageUrl",
				"email":    "$customer.email",
				"amount":   "$amount",
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	resp := make([]model.LatestInvoice, 0, 100)
	if err := cursor.All(ctx, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
