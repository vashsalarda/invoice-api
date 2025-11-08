package model

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type Revenue struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Month     string             `bson:"month" json:"month"`
	Year      string             `bson:"year" json:"year"`
	Revenue   float64            `bson:"revenue" json:"revenue"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

type RevenueDTO struct {
	ID      primitive.ObjectID `json:"id"`
	Month   string             `json:"month"`
	Year    string             `json:"year"`
	Revenue float64            `json:"revenue"`
}

type CreateRevenue struct {
	Month   string  `json:"month" validate:"required"`
	Year    string  `json:"year" validate:"required"`
	Revenue float64 `json:"revenue" validate:"required"`
}

type UpdateRevenue struct {
	Month   string  `json:"month"`
	Year    string  `json:"year"`
	Revenue float64 `json:"revenue"`
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
