package query

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/user/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DefaultQuery struct{}

func (c *DefaultQuery) CollectionName() string {
	return "users"
}

//go:generate mockgen -destination=../mocks/query/mock_pos_query.go -package=query shards.project-moonshot.com/ISAAC/go-inventory/features/pos/query PosQuery
type Query interface {
	GetAllByQuery() ([]model.UserDTO, error)
	GetByID(id string) (*model.UserDTO, error)
}

// GetAllUsers executes the get all users query
func (c *DefaultQuery) GetAllByQuery() ([]model.UserDTO, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{})
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

// GetUserByID executes the get user by ID query
func (c *DefaultQuery) GetByID(id string) (*model.UserDTO, error) {
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