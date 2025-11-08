package command

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/revenue/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultRevenueCommand struct{}

func (c *DefaultRevenueCommand) CollectionName() string {
	return "revenues"
}

type RevenueCommand interface {
	CreateItem(_val *model.CreateRevenue) (*mongo.InsertOneResult, error)
	UpdateItem(id string, _val *model.UpdateRevenue) (*mongo.UpdateResult, error)
	DeleteItem(id string) (*mongo.DeleteResult, error)
}

func (c *DefaultRevenueCommand) CreateItem(_val *model.CreateRevenue) (*mongo.InsertOneResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	doc := &model.Revenue{
		Month:     _val.Month,
		Year:      _val.Year,
		Revenue:   _val.Revenue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DefaultRevenueCommand) UpdateItem(id string, _val *model.UpdateRevenue) (*mongo.UpdateResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"month":     _val.Month,
			"year":      _val.Year,
			"revenue":   _val.Revenue,
			"updatedAt": time.Now(),
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

// DeleteRevenue executes the delete user command
func (c *DefaultRevenueCommand) DeleteItem(id string) (*mongo.DeleteResult, error) {
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
