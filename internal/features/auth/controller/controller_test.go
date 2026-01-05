package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	usermodel "invoice-api/internal/features/user/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthQuery struct{
    byEmail func(email string) (usermodel.User, error)
    byID func(id string) (*usermodel.UserDTO, error)
}

func (m *mockAuthQuery) GetItemByEmail(email string) (usermodel.User, error) {
    if m.byEmail == nil {
        return usermodel.User{}, nil
    }
    return m.byEmail(email)
}

func (m *mockAuthQuery) GetItemByID(id string) (*usermodel.UserDTO, error) {
    if m.byID == nil {
        return &usermodel.UserDTO{}, nil
    }
    return m.byID(id)
}

func (m *mockAuthQuery) GetItemsByQuery() ([]usermodel.UserDTO, error) {
    return []usermodel.UserDTO{}, nil
}

type mockAuthCommand struct{
    create func(u *usermodel.CreateUser) (*mongo.InsertOneResult, error)
}

func (m *mockAuthCommand) CreateUser(u *usermodel.CreateUser) (*mongo.InsertOneResult, error) {
    if m.create == nil {
        return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
    }
    return m.create(u)
}

func (m *mockAuthCommand) UpdateUser(id string, val *usermodel.UpdateUser) (*mongo.UpdateResult, error) {
    return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (m *mockAuthCommand) DeleteUser(id string) (*mongo.DeleteResult, error) {
    return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func TestGetUser_NotFoundAndSuccess(t *testing.T) {
    app := fiber.New()

    // not found: return zero user and nil error -> controller returns 404
    ctrl := &AuthController{Query: &mockAuthQuery{byID: func(id string) (*usermodel.UserDTO, error) {
        return nil, mongo.ErrNoDocuments
    }}}
    app.Get("/user/:id", ctrl.GetUser)
    r, _ := http.NewRequest("GET", "/user/abc", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 404 { t.Fatalf("expected 404 got %d", resp.StatusCode) }

    // success: return a user with ID
    id := primitive.NewObjectID()
    user := usermodel.UserDTO{ID: id, Email: "a@b.com"}
    ctrl2 := &AuthController{Query: &mockAuthQuery{byID: func(id string) (*usermodel.UserDTO, error) {
        return &user, nil
    }}}
    app2 := fiber.New()
    app2.Get("/user/:id", ctrl2.GetUser)
    r2, _ := http.NewRequest("GET", "/user/507f1f77bcf86cd799439011", nil)
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp2.StatusCode) }
}

func TestSignInUser_BodyParserAndAuth(t *testing.T) {
    app := fiber.New()

    // Body parser error -> invalid JSON
    ctrl := &AuthController{Query: &mockAuthQuery{}}
    app.Post("/signin", ctrl.SignInUser)
    r, _ := http.NewRequest("POST", "/signin", bytes.NewReader([]byte("invalid")))
    r.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 422 { t.Fatalf("expected 422 got %d", resp.StatusCode) }

    // Wrong password
    pw, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
    ctrl2 := &AuthController{Query: &mockAuthQuery{byEmail: func(email string) (usermodel.User, error) {
        return usermodel.User{Email: email, Password: string(pw)}, nil
    }}}
    app2 := fiber.New()
    app2.Post("/signin", ctrl2.SignInUser)
    body := bytes.NewReader([]byte(`{"email":"a@b.com","password":"wrongpass"}`))
    r2, _ := http.NewRequest("POST", "/signin", body)
    r2.Header.Set("Content-Type", "application/json")
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp2.StatusCode) }

    // Success: correct password -> returns token
    ctrl3 := &AuthController{Query: &mockAuthQuery{byEmail: func(email string) (usermodel.User, error) {
        return usermodel.User{ID: primitive.NewObjectID(), Email: email, Password: string(pw)}, nil
    }}}
    os.Setenv("JWT_SECRET", "testsecret")
    app3 := fiber.New()
    app3.Post("/signin", ctrl3.SignInUser)
    body3 := bytes.NewReader([]byte(`{"email":"a@b.com","password":"correctpass"}`))
    r3, _ := http.NewRequest("POST", "/signin", body3)
    r3.Header.Set("Content-Type", "application/json")
    resp3, err := app3.Test(r3)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp3.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp3.StatusCode) }
    var got map[string]interface{}
    json.NewDecoder(resp3.Body).Decode(&got)
    if got["token"] == "" || got["token"] == nil {
        t.Fatalf("expected token in response")
    }
}

func TestSignUpUser_ValidationAndCreate(t *testing.T) {
    app := fiber.New()

    // invalid json -> BodyParser error
    ctrl := &AuthController{Query: &mockAuthQuery{}, Command: &mockAuthCommand{}}
    app.Post("/signup", ctrl.SignUpUser)
    r, _ := http.NewRequest("POST", "/signup", bytes.NewReader([]byte("invalid")))
    r.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp.StatusCode) }

    // validation error: empty payload
    app2 := fiber.New()
    app2.Post("/signup", ctrl.SignUpUser)
    r2, _ := http.NewRequest("POST", "/signup", bytes.NewReader([]byte(`{}`)))
    r2.Header.Set("Content-Type", "application/json")
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp2.StatusCode) }
}

func TestLogoutUser(t *testing.T) {
    app := fiber.New()
    ctrl := &AuthController{}
    app.Post("/logout", ctrl.LogoutUser)
    r, _ := http.NewRequest("POST", "/logout", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp.StatusCode) }
}
