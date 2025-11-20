package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/invoice/model"
	"math"
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
	GetItemsByQuery(keyword string, status string, size int64, page int64) (*model.InvoicePage, error)
	GetItemByID(id string) (*model.InvoiceDTO, error)
	GetLatestInvoices() ([]model.LatestInvoice, error)
	GetTotalItemsByQuery(keyword string, status string) (int64, error)
}

func (c *DefaultInvoiceQuery) GetItemsByQuery(keyword string, status string, size int64, page int64) (*model.InvoicePage, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	var filter = bson.M{}
	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"middleName": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	if status != "" {
		filter["status"] = status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int64(math.Ceil(float64(totalItems) / float64(page)))

	items := make([]*model.InvoiceDTO, 0, 100)
	opts := new(options.FindOptions)
	if size != 0 {
		if page == 0 {
			page = 1
		}
		opts.SetSkip(int64((page - 1) * size))
		opts.SetLimit(int64(size))
		opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	}

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	resp := &model.InvoicePage{
		TotalRows:  int64(len(items)),
		TotalPages: totalPages,
		PageNumber: page,
		PageSize:   page,
		Data:       items,
	}

	return resp, nil
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

func (c *DefaultInvoiceQuery) GetTotalItemsByQuery(keyword string, status string) (int64, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	var filter = bson.M{}

	if keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"customerId": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"middleName": bson.M{"$regex": keyword, "$options": "i"}},
			bson.M{"email": bson.M{"$regex": keyword, "$options": "i"}},
		}
	}

	if status != "" {
		filter["status"] = status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return totalItems, nil
}
