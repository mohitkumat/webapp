package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"api/internal/handlers"
)

func TestListItems_WhenEmpty_ReturnsEmptyArrayNotNull(t *testing.T) {
	conn := setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()

	handlers.ListItems(conn)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "[]\n" {
		t.Fatalf("expected empty JSON array, got %q", rec.Body.String())
	}
}

func TestCreateItem_WithValidName_PersistsAndReturnsIt(t *testing.T) {
	conn := setupTestDB(t)

	body, _ := json.Marshal(map[string]string{"name": "first widget"})
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.CreateItem(conn)(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	var created handlers.Item
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created.Name != "first widget" || created.ID == 0 {
		t.Fatalf("unexpected created item: %+v", created)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/items", nil)
	listRec := httptest.NewRecorder()
	handlers.ListItems(conn)(listRec, listReq)

	var items []handlers.Item
	if err := json.NewDecoder(listRec.Body).Decode(&items); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(items) != 1 || items[0].Name != "first widget" {
		t.Fatalf("expected one persisted item, got %+v", items)
	}
}

func TestCreateItem_WithEmptyName_Returns422(t *testing.T) {
	conn := setupTestDB(t)

	body, _ := json.Marshal(map[string]string{"name": ""})
	req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handlers.CreateItem(conn)(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}
