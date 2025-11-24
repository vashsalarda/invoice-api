package command

import (
	"context"
	"errors"
	"invoice-api/internal/database"
	customer_model "invoice-api/internal/features/customer/model"
	"invoice-api/internal/features/invoice/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultInvoiceCommand struct{}

func (c *DefaultInvoiceCommand) CollectionName() string {
	return "invoices"
}

type InvoiceCommand interface {
	CreateItem(_val *model.CreateInvoice) (*mongo.InsertOneResult, error)
	UpdateItem(id string, _val *model.UpdateInvoice) (*mongo.UpdateResult, error)
	DeleteItem(id string) (*mongo.DeleteResult, error)
}

func (c *DefaultInvoiceCommand) CreateItem(_val *model.CreateInvoice) (*mongo.InsertOneResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())
	customerCollection := db.Collection("customers")

	customerID, err := primitive.ObjectIDFromHex(_val.CustomerID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	customer := customer_model.CustomerDTOMin{}
	err = customerCollection.FindOne(ctx, bson.M{"_id": customerID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	doc := &model.Invoice{
		CustomerID: customerID,
		Customer:   customer,
		Status:     _val.Status,
		Amount:     _val.Amount,
		Date:       _val.Date,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DefaultInvoiceCommand) UpdateItem(id string, _val *model.UpdateInvoice) (*mongo.UpdateResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	customerID, err := primitive.ObjectIDFromHex(_val.CustomerID)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"customerID": customerID,
			"status":     _val.Status,
			"amount":     _val.Amount,
			"date":       _val.Date,
			"updatedAt":  time.Now(),
		},
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return result, nil
}

// DeleteInvoice executes the delete user command
func (c *DefaultInvoiceCommand) DeleteItem(id string) (*mongo.DeleteResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return nil, err
	}

	if res.DeletedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return res, nil
}
