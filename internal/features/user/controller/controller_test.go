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

	modelpkg "invoice-api/internal/features/user/model"
)

type mockCommand struct {
	createRes *mongo.InsertOneResult
	createErr error
	updateRes *mongo.UpdateResult
	updateErr error
	deleteRes *mongo.DeleteResult
	deleteErr error
}

func (m *mockCommand) CreateUser(_val *modelpkg.CreateUser) (*mongo.InsertOneResult, error) {
	return m.createRes, m.createErr
}

func (m *mockCommand) UpdateUser(id string, _val *modelpkg.UpdateUser) (*mongo.UpdateResult, error) {
	return m.updateRes, m.updateErr
}

func (m *mockCommand) DeleteUser(id string) (*mongo.DeleteResult, error) {
	return m.deleteRes, m.deleteErr
}

type mockQuery struct {
	itemsRes []modelpkg.UserDTO
	itemsErr error
	itemRes  *modelpkg.UserDTO
	itemErr  error
	byEmail  modelpkg.User
	byEmailErr error
}

func (m *mockQuery) GetItemsByQuery() ([]modelpkg.UserDTO, error) {
	return m.itemsRes, m.itemsErr
}

func (m *mockQuery) GetItemByID(id string) (*modelpkg.UserDTO, error) {
	return m.itemRes, m.itemErr
}

func (m *mockQuery) GetItemByEmail(email string) (modelpkg.User, error) {
	return m.byEmail, m.byEmailErr
}

func TestCreateUser_SuccessAndBadBodyAndValidationAndEmailTaken(t *testing.T) {
	app := fiber.New()

	// success case: email not taken
	mockQ := &mockQuery{byEmail: modelpkg.User{ID: primitive.NilObjectID}}
	mockC := &mockCommand{createRes: &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, createErr: nil}
	ctrl := &UserController{Command: mockC, Query: mockQ}
	app.Post("/", ctrl.CreateUser)

	payload := map[string]string{"firstName": "John", "lastName": "Doe", "email": "john@example.com", "password": "secret"}
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
	payload3 := map[string]string{"firstName": "John", "lastName": "", "email": "", "password": ""}
	b3, _ := json.Marshal(payload3)
	req3 := httptest.NewRequest("POST", "/", bytes.NewReader(b3))
	req3.Header.Set("Content-Type", "application/json")
	resp3, err := app.Test(req3)
	require.NoError(t, err)
	require.Equal(t, 400, resp3.StatusCode)

	// email taken
	takenQ := &mockQuery{byEmail: modelpkg.User{ID: primitive.NewObjectID()}}
	takenCtrl := &UserController{Command: mockC, Query: takenQ}
	app2 := fiber.New()
	app2.Post("/", takenCtrl.CreateUser)
	b4, _ := json.Marshal(payload)
	req4 := httptest.NewRequest("POST", "/", bytes.NewReader(b4))
	req4.Header.Set("Content-Type", "application/json")
	resp4, err := app2.Test(req4)
	require.NoError(t, err)
	require.Equal(t, 409, resp4.StatusCode)
}

func TestGetAllUsers_ErrorAndSuccess(t *testing.T) {
	app := fiber.New()

	// error case
	errMock := &mockQuery{itemsRes: nil, itemsErr: mongo.ErrNoDocuments}
	ctrl := &UserController{Query: errMock}
	app.Get("/", ctrl.GetAllUsers)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)

	// success case
	user := modelpkg.UserDTO{ID: primitive.NewObjectID(), FirstName: "A", LastName: "B", Email: "a@b.com"}
	succMock := &mockQuery{itemsRes: []modelpkg.UserDTO{user}, itemsErr: nil}
	ctrl2 := &UserController{Query: succMock}
	app2 := fiber.New()
	app2.Get("/", ctrl2.GetAllUsers)

	req2 := httptest.NewRequest("GET", "/", nil)
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
}

func TestGetUserByID_NotFoundAndSuccess(t *testing.T) {
	app := fiber.New()

	// not found
	notFoundMock := &mockQuery{itemRes: nil, itemErr: mongo.ErrNoDocuments}
	ctrl := &UserController{Query: notFoundMock}
	app.Get("/:id", ctrl.GetUserByID)

	req := httptest.NewRequest("GET", "/someid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// success
	user := &modelpkg.UserDTO{ID: primitive.NewObjectID(), FirstName: "A", LastName: "B", Email: "a@b.com"}
	succMock := &mockQuery{itemRes: user, itemErr: nil}
	ctrl2 := &UserController{Query: succMock}
	app2 := fiber.New()
	app2.Get("/:id", ctrl2.GetUserByID)

	req2 := httptest.NewRequest("GET", "/someid", nil)
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
}

func TestUpdateUser_NotFoundAndSuccess(t *testing.T) {
	// not found (err == mongo.ErrNoDocuments)
	app := fiber.New()
	notFoundMock := &mockCommand{updateRes: nil, updateErr: mongo.ErrNoDocuments}
	ctrl := &UserController{Command: notFoundMock}
	app.Put("/:id", ctrl.UpdateUser)

	payload := map[string]string{"firstName": "New Name"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("PUT", "/someid", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// update success with ModifiedCount 1
	succMock := &mockCommand{updateRes: &mongo.UpdateResult{ModifiedCount: 1}, updateErr: nil}
	ctrl2 := &UserController{Command: succMock}
	app2 := fiber.New()
	app2.Put("/:id", ctrl2.UpdateUser)

	req2 := httptest.NewRequest("PUT", "/someid", bytes.NewReader(b))
	req2.Header.Set("Content-Type", "application/json")
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
	bb, _ := io.ReadAll(resp2.Body)
	var body map[string]string
	err = json.Unmarshal(bb, &body)
	require.NoError(t, err)
	require.Equal(t, "User updated successfully", body["message"])
}

func TestDeleteUser_NotFoundAndSuccess(t *testing.T) {
	app := fiber.New()

	// Not found case
	notFoundMock := &mockCommand{deleteRes: nil, deleteErr: mongo.ErrNoDocuments}
	ctrl := &UserController{Command: notFoundMock}
	app.Delete("/:id", ctrl.DeleteUser)

	req := httptest.NewRequest("DELETE", "/someid", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 404, resp.StatusCode)

	// success case
	succMock := &mockCommand{deleteRes: &mongo.DeleteResult{DeletedCount: 1}, deleteErr: nil}
	ctrl2 := &UserController{Command: succMock}
	app2 := fiber.New()
	app2.Delete("/:id", ctrl2.DeleteUser)

	req2 := httptest.NewRequest("DELETE", "/someid", nil)
	resp2, err := app2.Test(req2)
	require.NoError(t, err)
	require.Equal(t, 200, resp2.StatusCode)
	b, _ := io.ReadAll(resp2.Body)
	var body map[string]string
	err = json.Unmarshal(b, &body)
	require.NoError(t, err)
	require.Equal(t, "User deleted successfully", body["message"])
}
