package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"invoice-api/internal/features/invoice/model"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mockQuery struct{
    getByID func(id string) (*model.InvoiceDTO, error)
    getTotal func(keyword string, status string) (int64, error)
}

func (m *mockQuery) GetItemsByQuery(keyword string, status string, size int64, page int64) (*model.InvoicePage, error) {
    return &model.InvoicePage{}, nil
}
func (m *mockQuery) GetItemByID(id string) (*model.InvoiceDTO, error) {
    return m.getByID(id)
}
func (m *mockQuery) GetLatestInvoices() ([]model.LatestInvoice, error) {
    return nil, nil
}
func (m *mockQuery) GetTotalItemsByQuery(keyword string, status string) (int64, error) {
    return m.getTotal(keyword, status)
}
func (m *mockQuery) GetCustomersInvoices(keyword string) ([]model.InvoiceCustomers, error) {
    return nil, nil
}

type mockCommand struct{
	createRes *mongo.InsertOneResult
	createErr error
    update func(id string, val *model.UpdateInvoice) (*mongo.UpdateResult, error)
    del func(id string) (*mongo.DeleteResult, error)
}

func (m *mockCommand) CreateCustomer(_val *model.CreateInvoice) (*mongo.InsertOneResult, error) {
	return m.createRes, m.createErr
}

func (m *mockCommand) CreateItem(_val *model.CreateInvoice) (*mongo.InsertOneResult, error) {
    return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
}
func (m *mockCommand) UpdateItem(id string, _val *model.UpdateInvoice) (*mongo.UpdateResult, error) {
    return m.update(id, _val)
}
func (m *mockCommand) DeleteItem(id string) (*mongo.DeleteResult, error) {
    return m.del(id)
}

func TestCreateCustomer_Success(t *testing.T) {
	app := fiber.New()

	// success case
	mock := &mockCommand{createRes: &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, createErr: nil}
	ctrl := &InvoiceController{Command: mock}
	app.Post("/invoices", ctrl.CreateInvoice)

	payload := model.CreateInvoice{
		CustomerID: "507f1f77bcf86cd799439011",
		Amount:     100.50,
		Date:       "2024-06-01",
		Status:     "pending",
	}
	
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/invoices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, 201, resp.StatusCode)
}

func TestCreateInvoice_ValidationError(t *testing.T) {
    app := fiber.New()
    ctrl := &InvoiceController{}
    app.Post("/invoices", ctrl.CreateInvoice)

    // send empty object -> validation should fail (required fields missing)
    req := bytes.NewReader([]byte(`{}`))
    r, _ := http.NewRequest("POST", "/invoices", req)
    r.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(r)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    if resp.StatusCode != 400 {
        t.Fatalf("expected 400 got %d", resp.StatusCode)
    }
}

func TestCreateInvoice_InvalidPayloadError(t *testing.T) {
    app := fiber.New()
    ctrl := &InvoiceController{}
    app.Post("/invoices", ctrl.CreateInvoice)

    req := bytes.NewReader([]byte(`invalid-json`))
    r, _ := http.NewRequest("POST", "/invoices", req)
    r.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(r)
    if err != nil {
        t.Fatalf("request failed: %v", err)
    }
    if resp.StatusCode != 400 {
        t.Fatalf("expected 400 got %d", resp.StatusCode)
    }
}

func TestGetInvoiceByID_NotFoundAndSuccess(t *testing.T) {
    app := fiber.New()
    // not found
    ctrl := &InvoiceController{Query: &mockQuery{getByID: func(id string) (*model.InvoiceDTO, error) {
        return nil, mongo.ErrNoDocuments
    }}}
    app.Get("/invoices/:id", ctrl.GetInvoiceByID)

    r, _ := http.NewRequest("GET", "/invoices/abc", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 404 { t.Fatalf("expected 404 got %d", resp.StatusCode) }

    // success
    id := primitive.NewObjectID()
    dto := &model.InvoiceDTO{ID: id}
    ctrl2 := &InvoiceController{Query: &mockQuery{getByID: func(id string) (*model.InvoiceDTO, error) {
        return dto, nil
    }}}
    app2 := fiber.New()
    app2.Get("/invoices/:id", ctrl2.GetInvoiceByID)
    r2, _ := http.NewRequest("GET", "/invoices/abc", nil)
    resp2, err := app2.Test(r2)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp2.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp2.StatusCode) }
}

func TestUpdateInvoice_FailedToUpdate(t *testing.T) {
    app := fiber.New()
    ctrl := &InvoiceController{Command: &mockCommand{update: func(id string, val *model.UpdateInvoice) (*mongo.UpdateResult, error) {
        return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 0}, nil
    }}}
    app.Put("/invoices/:id", ctrl.UpdateInvoice)

    body := bytes.NewReader([]byte(`{"amount":100}`))
    r, _ := http.NewRequest("PUT", "/invoices/507f1f77bcf86cd799439011", body)
    r.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 404 { t.Fatalf("expected 404 got %d", resp.StatusCode) }
}

func TestDeleteInvoice_Success(t *testing.T) {
    app := fiber.New()
    ctrl := &InvoiceController{Command: &mockCommand{del: func(id string) (*mongo.DeleteResult, error) {
        return &mongo.DeleteResult{DeletedCount: 1}, nil
    }}}
    app.Delete("/invoices/:id", ctrl.DeleteInvoice)

    r, _ := http.NewRequest("DELETE", "/invoices/507f1f77bcf86cd799439011", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 200 { t.Fatalf("expected 200 got %d", resp.StatusCode) }

    // assert message in response
    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    if body["message"] != "Invoice deleted successfully" {
        t.Fatalf("unexpected message: %v", body)
    }
}

func TestGetTotalInvoices_InvalidStatus(t *testing.T) {
    app := fiber.New()
    ctrl := &InvoiceController{Query: &mockQuery{getTotal: func(k, s string) (int64, error) { return 0, nil }}}
    app.Get("/invoices/total", ctrl.GetTotalInvoices)

    r, _ := http.NewRequest("GET", "/invoices/total?status=badstatus", nil)
    resp, err := app.Test(r)
    if err != nil { t.Fatalf("request failed: %v", err) }
    if resp.StatusCode != 400 { t.Fatalf("expected 400 got %d", resp.StatusCode) }
}
