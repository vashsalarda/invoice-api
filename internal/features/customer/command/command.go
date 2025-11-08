package command

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/customer/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultCommand struct{}

func (c *DefaultCommand) CollectionName() string {
	return "customers"
}

type Command interface {
	CreateCustomer(_val *model.CreateCustomer) (*mongo.InsertOneResult, error)
	UpdateCustomer(id string, _val *model.UpdateCustomer) (*mongo.UpdateResult, error)
	DeleteCustomer(id string) (*mongo.DeleteResult, error)
}

func (c *DefaultCommand) CreateCustomer(_val *model.CreateCustomer) (*mongo.InsertOneResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	imageURL := _val.ImageURL
	if _val.ImageURL == "" {
		imageURL = "https://placehold.co/250/93C5fd/fff/png?text=" + _val.Name[:1]
	}

	customer := &model.Customer{
		Name:  _val.Name,
		Email:      _val.Email,
		ImageURL:   imageURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DefaultCommand) UpdateCustomer(id string, _val *model.UpdateCustomer) (*mongo.UpdateResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":  _val.Name,
			"imageUrl":   _val.ImageURL,
			"updatedAt":  time.Now(),
		},
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return nil, err
	}

	if res.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return res, nil
}

func (c *DefaultCommand) DeleteCustomer(id string) (*mongo.DeleteResult, error) {
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