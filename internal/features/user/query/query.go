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
	GetAllByQuery() ([]model.User, error)
	GetByID(id string) (*model.User, error)
}

// GetAllUsers executes the get all users query
func (c *DefaultQuery) GetAllByQuery() ([]model.User, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	users := make([]model.User, 0, 100) 
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetUserByID executes the get user by ID query
func (c *DefaultQuery) GetByID(id string) (*model.User, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())
	
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user model.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}