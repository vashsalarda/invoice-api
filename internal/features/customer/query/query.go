package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/customer/model"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultQuery struct{}

func (c *DefaultQuery) CollectionName() string {
	return "customers"
}

type Query interface {
	GetItemsByQuery(keyword string, size int64, page int64) (*model.CustomerPage, error)
	GetItemByID(id string) (*model.CustomerDTO, error)
	GetByEmail(email string) (model.CustomerDTO, error)
	GetTotalItemsByQuery(keyword string) (int64, error)
	GetItemsWithTotalByQuery(keyword string, size int64, page int64) (*model.CustomerWithTotalPage, error)
}

func (c *DefaultQuery) GetItemsByQuery(keyword string, size int64, page int64) (*model.CustomerPage, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int64(math.Ceil(float64(totalItems) / float64(page)))

	items := make([]*model.CustomerDTO, 0, 100)
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

	resp := &model.CustomerPage{
		TotalRows:  int64(len(items)),
		TotalPages: totalPages,
		PageNumber: page,
		PageSize:   page,
		Data:       items,
	}

	return resp, nil
}

func (c *DefaultQuery) GetItemByID(id string) (*model.CustomerDTO, error) {
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

func (c *DefaultQuery) GetByEmail(email string) (model.CustomerDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	var user model.CustomerDTO
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (c *DefaultQuery) GetTotalItemsByQuery(keyword string) (int64, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return totalItems, nil
}

func (c *DefaultQuery) GetItemsWithTotalByQuery(keyword string, size int64, page int64) (*model.CustomerWithTotalPage, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	totalItems, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	totalPages := int64(math.Ceil(float64(totalItems) / float64(page)))

	items := make([]*model.CustomerWithTotalDTO, 0, 100)
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

	for i, item := range items {
		if !item.ID.IsZero() {
			invoiceCollection := db.Collection("invoices")
			totalInvoices, err := invoiceCollection.CountDocuments(ctx, bson.M{"customer._id": item.ID})
			if err != nil {
				totalInvoices = 0
			}
			totalPaid, err := invoiceCollection.CountDocuments(ctx, bson.M{"customer._id": item.ID, "status": "paid"})
			if err != nil {
				totalPaid = 0
			}
			totalPending, err := invoiceCollection.CountDocuments(ctx, bson.M{"customer._id": item.ID, "status": bson.M{"$ne": "paid"}})
			if err != nil {
				totalPaid = 0
			}

			items[i].TotalInvoices = totalInvoices
			items[i].TotalPaid = totalPaid
			items[i].TotalPending = totalPending
		}
	}

	resp := &model.CustomerWithTotalPage{
		TotalRows:  int64(len(items)),
		TotalPages: totalPages,
		PageNumber: page,
		PageSize:   page,
		Data:       items,
	}

	return resp, nil
}
