package model

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type Customer struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName  string             `bson:"firstName" json:"firstName"`
	LastName   string             `bson:"lastName" json:"lastName"`
	MiddleName string             `bson:"middleName" json:"middleName"`
	Email      string             `bson:"email" json:"email"`
	ImageURL   string             `bson:"imageUrl" json:"imageUrl"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type CustomerDTO struct {
	ID         primitive.ObjectID `json:"id"`
	FirstName  string             `json:"firstName"`
	LastName   string             `json:"lastName"`
	MiddleName string             `json:"middleName"`
	Email      string             `json:"email"`
	ImageURL   string             `json:"imageUrl"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt,omitzero"`
}

type CreateCustomer struct {
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
	MiddleName string `json:"middleName"`
	Email      string `json:"email" validate:"required"`
	ImageURL   string `json:"imageUrl"`
}

type UpdateCustomer struct {
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
	ImageURL   string `json:"imageUrl"`
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []ErrorResponse {
	var errors []ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}
	return errors
}
