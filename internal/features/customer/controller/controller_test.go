package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	modelpkg "invoice-api/internal/features/customer/model"
)

type commandMock struct {
	createRes *mongo.InsertOneResult
	createErr error
	updateRes *mongo.UpdateResult
	updateErr error
	deleteRes *mongo.DeleteResult
	deleteErr error
}

func (m *commandMock) CreateCustomer(_val *modelpkg.CreateCustomer) (*mongo.InsertOneResult, error) {
	return m.createRes, m.createErr
}

func (m *commandMock) UpdateCustomer(id string, _val *modelpkg.UpdateCustomer) (*mongo.UpdateResult, error) {
	return m.updateRes, m.updateErr
}

func (m *commandMock) DeleteCustomer(id string) (*mongo.DeleteResult, error) {
	return m.deleteRes, m.deleteErr
}

func TestCreateCustomer_SuccessAndBadBody(t *testing.T) {
	app := fiber.New()

	// success case
	mock := &commandMock{createRes: &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, createErr: nil}
	ctrl := &CustomerController{Command: mock}
	app.Post("/", ctrl.CreateCustomer)

	payload := map[string]string{"name": "John", "email": "john@example.com"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)

	// bad body
	req2 := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("{badjson")))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 400, resp2.StatusCode)

	// validation error (missing fields)
	payload3 := map[string]string{"name": "John", "email": ""}
	b3, _ := json.Marshal(payload3)
	req3 := httptest.NewRequest("POST", "/", bytes.NewReader(b3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := app.Test(req3)
	require.NoError(t, err)
	require.Equal(t, 400, resp3.StatusCode)
}

func TestDeleteCustomer_NotFoundAndSuccess(t *testing.T) {
	app := fiber.New()

	// Not found case
	notFoundMock := &commandMock{deleteRes: nil, deleteErr: mongo.ErrNoDocuments}
	ctrl := &CustomerController{Command: notFoundMock}
	app.Delete("/:id", ctrl.DeleteCustomer)

	req := httptest.NewRequest("DELETE", "/someid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// success case
	succMock := &commandMock{deleteRes: &mongo.DeleteResult{DeletedCount: 1}, deleteErr: nil}
	ctrl2 := &CustomerController{Command: succMock}
	app2 := fiber.New()
	app2.Delete("/:id", ctrl2.DeleteCustomer)

	req2 := httptest.NewRequest("DELETE", "/someid", nil)
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
	b, _ := io.ReadAll(resp2.Body)
	var body map[string]string
	err = json.Unmarshal(b, &body)
	require.NoError(t, err)
	require.Equal(t, "Customer deleted successfully", body["message"])
}

func TestUpdateCustomer_NotFoundAndSuccess(t *testing.T) {
	// not found (err == mongo.ErrNoDocuments)
	app := fiber.New()
	notFoundMock := &commandMock{updateRes: nil, updateErr: mongo.ErrNoDocuments}
	ctrl := &CustomerController{Command: notFoundMock}
	app.Put("/:id", ctrl.UpdateCustomer)

	payload := map[string]string{"name": "New Name"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/someid", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// update success with ModifiedCount 1
	succMock := &commandMock{updateRes: &mongo.UpdateResult{ModifiedCount: 1}, updateErr: nil}
	ctrl2 := &CustomerController{Command: succMock}
	app2 := fiber.New()
	app2.Put("/:id", ctrl2.UpdateCustomer)

	req2 := httptest.NewRequest("PUT", "/someid", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
	bb, _ := io.ReadAll(resp2.Body)
	var body map[string]string
	err = json.Unmarshal(bb, &body)
	require.NoError(t, err)
	require.Equal(t, "Customer updated successfully", body["message"])
}
