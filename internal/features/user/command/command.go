package command

import (
	"context"
	"invoice-api/internal/database"
	"invoice-api/internal/features/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type DefaultCommand struct{}
func (c *DefaultCommand) CollectionName() string {
	return "users"
}


type Command interface {
	CreateUser(_val *model.CreateUser) (*mongo.InsertOneResult, error)
	UpdateUser(id string, _val *model.UpdateUser) (*mongo.UpdateResult, error)
	DeleteUser(id string) (*mongo.DeleteResult, error)
}

func (c *DefaultCommand) CreateUser(_val *model.CreateUser) (*mongo.InsertOneResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(_val.Password), bcrypt.DefaultCost)
	user := &model.User{
		Name:     _val.Name,
		Email:    _val.Email,
		Password: string(hashedPassword),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *DefaultCommand) UpdateUser(id string, val *model.UpdateUser) (*mongo.UpdateResult, error) {
	db := database.GetDatabase()
	collection := db.Collection(c.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":     val.Name,
			"email":    val.Email,
			"password": val.Password,
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

// DeleteUser executes the delete user command
func (c *DefaultCommand) DeleteUser(id string) (*mongo.DeleteResult, error) {
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
