package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0501jp/shinjuku-lunch/api/db"
)

func GetAreas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.QueryContext(r.Context(), "SELECT id, name FROM areas ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type area struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var areas []area
	for rows.Next() {
		var a area
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		areas = append(areas, a)
	}
	if areas == nil {
		areas = []area{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}
