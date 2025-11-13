package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DefaultQuery struct{}

func (c *DefaultQuery) CollectionName() string {
	return "users"
}

//go:generate mockgen -destination=../mocks/query/mock_pos_query.go -package=query shards.project-moonshot.com/ISAAC/go-inventory/features/pos/query PosQuery
type Query interface {
	GetItemsByQuery() ([]model.UserDTO, error)
	GetItemByID(id string) (*model.UserDTO, error)
	GetItemByEmail(email string) (model.User, error)
}

func (c *DefaultQuery) GetItemsByQuery() ([]model.UserDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	items := make([]model.UserDTO, 0, 100) 
	if err := cursor.All(context.TODO(), &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (c *DefaultQuery) GetItemByID(id string) (*model.UserDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var items model.UserDTO
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&items)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &items, nil
}

func (c *DefaultQuery) GetItemByEmail(email string) (model.User, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	var user model.User
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}