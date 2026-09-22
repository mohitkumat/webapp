package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type createItemRequest struct {
	Name string `json:"name"`
}

// ListItems returns all items. Demonstrates the empty-state contract: an empty
// array, not null, and never a bare 200 with no body.
func ListItems(conn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := conn.QueryContext(r.Context(), `SELECT id, name FROM items ORDER BY id`)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not list items"})
			return
		}
		defer rows.Close()

		items := []Item{}
		for rows.Next() {
			var it Item
			if err := rows.Scan(&it.ID, &it.Name); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not read item"})
				return
			}
			items = append(items, it)
		}

		writeJSON(w, http.StatusOK, items)
	}
}

// CreateItem demonstrates request validation and a structured error response.
func CreateItem(conn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
			return
		}
		if req.Name == "" {
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "name is required"})
			return
		}

		var it Item
		err := conn.QueryRowContext(
			r.Context(),
			`INSERT INTO items (name) VALUES ($1) RETURNING id, name`,
			req.Name,
		).Scan(&it.ID, &it.Name)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create item"})
			return
		}

		writeJSON(w, http.StatusCreated, it)
	}
}
