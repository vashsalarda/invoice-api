package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"invoice-api/internal/features/revenue/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockQuery struct{
    getAll func() ([]model.RevenueDTO, error)
    getByID func(id string) (*model.RevenueDTO, error)
}

func (m *mockQuery) GetItemsByQuery() ([]model.RevenueDTO, error) {
    if m.getAll == nil { return []model.RevenueDTO{}, nil }
    return m.getAll()
}
func (m *mockQuery) GetItemByID(id string) (*model.RevenueDTO, error) {
    if m.getByID == nil { return &model.RevenueDTO{}, nil }
    return m.getByID(id)
}

type mockCommand struct{
    create func(val *model.CreateRevenue) (*mongo.InsertOneResult, error)
	createRes *mongo.InsertOneResult
	createErr error
    update func(id string, val *model.UpdateRevenue) (*mongo.UpdateResult, error)
    del func(id string) (*mongo.DeleteResult, error)
}

func (m *mockCommand) CreateItem(val *model.CreateRevenue) (*mongo.InsertOneResult, error) {
    if m.create == nil { return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil }
    return m.create(val)
}
func (m *mockCommand) UpdateItem(id string, val *model.UpdateRevenue) (*mongo.UpdateResult, error) {
    if m.update == nil { return &mongo.UpdateResult{MatchedCount:1, ModifiedCount:1}, nil }
    return m.update(id, val)
}
func (m *mockCommand) DeleteItem(id string) (*mongo.DeleteResult, error) {
    if m.del == nil { return &mongo.DeleteResult{DeletedCount:1}, nil }
    return m.del(id)
}

func TestCreateCustomer_Success(t *testing.T) {
	app := fiber.New()

	// success case
	mock := &mockCommand{createRes: &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, createErr: nil}
	ctrl := &RevenueController{Command: mock}
	app.Post("/revenues", ctrl.CreateRevenue)

	payload := model.CreateRevenue{
		Month: "01",
		Year:  "2025",
		Revenue: 10000,
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/revenues", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
}

func TestCreateRevenue_ValidationAndBodyParser(t *testing.T) {
    app := fiber.New()
    ctrl := &RevenueController{}
    app.Post("/revenues", ctrl.CreateRevenue)

    // invalid JSON -> body parser error
    r, _ := http.NewRequest("POST", "/revenues", bytes.NewReader([]byte("invalid")))
    r.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp.StatusCode) }

    // validation error: missing required fields
    app2 := fiber.New()
    app2.Post("/revenues", ctrl.CreateRevenue)
    r2, _ := http.NewRequest("POST", "/revenues", bytes.NewReader([]byte(`{}`)))
    r2.Header.Set("Content-Type", "application/json")
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp2.StatusCode) }
}

func TestGetAllRevenues_Success(t *testing.T) {
    app := fiber.New()
    ctrl := &RevenueController{Query: &mockQuery{getAll: func() ([]model.RevenueDTO, error) {
        return []model.RevenueDTO{{Month: "Jan"}}, nil
    }}}
    app.Get("/revenues", ctrl.GetAllRevenues)
    r, _ := http.NewRequest("GET", "/revenues", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp.StatusCode) }
}

func TestGetRevenueByID_NotFoundAndSuccess(t *testing.T) {
    app := fiber.New()
    // not found -> simulate mongo.ErrNoDocuments via returning nil and error
    ctrl := &RevenueController{Query: &mockQuery{getByID: func(id string) (*model.RevenueDTO, error) {
        return nil, mongo.ErrNoDocuments
    }}}
    app.Get("/revenues/:id", ctrl.GetRevenueByID)
    r, _ := http.NewRequest("GET", "/revenues/abc", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 404 { t.Fatalf("expected 404 got %d", resp.StatusCode) }

    // success
    dto := &model.RevenueDTO{Month: "Feb"}
    ctrl2 := &RevenueController{Query: &mockQuery{getByID: func(id string) (*model.RevenueDTO, error) {
        return dto, nil
    }}}
    app2 := fiber.New()
    app2.Get("/revenues/:id", ctrl2.GetRevenueByID)
    r2, _ := http.NewRequest("GET", "/revenues/507f1f77bcf86cd799439011", nil)
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp2.StatusCode) }
}

func TestUpdateRevenue_ModifiedCountZero(t *testing.T) {
    app := fiber.New()
    ctrl := &RevenueController{Command: &mockCommand{update: func(id string, val *model.UpdateRevenue) (*mongo.UpdateResult, error) {
        return &mongo.UpdateResult{MatchedCount:1, ModifiedCount:0}, nil
    }}}
    app.Put("/revenues/:id", ctrl.UpdateRevenue)
    body := bytes.NewReader([]byte(`{"revenue":123}`))
    r, _ := http.NewRequest("PUT", "/revenues/507f1f77bcf86cd799439011", body)
    r.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 404 { t.Fatalf("expected 404 got %d", resp.StatusCode) }
}

func TestDeleteRevenue_Success(t *testing.T) {
    app := fiber.New()
    ctrl := &RevenueController{Command: &mockCommand{del: func(id string) (*mongo.DeleteResult, error) {
        return &mongo.DeleteResult{DeletedCount:1}, nil
    }}}
    app.Delete("/revenues/:id", ctrl.DeleteRevenue)
    r, _ := http.NewRequest("DELETE", "/revenues/507f1f77bcf86cd799439011", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp.StatusCode) }
    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    if body["message"] != "Revenue deleted successfully" { t.Fatalf("unexpected message: %v", body) }
}
