package handler

import (
	"encoding/json"
	"net/http"

	"github.com/0501jp/shinjuku-lunch/api/db"
)

func GetGenres(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.QueryContext(r.Context(), "SELECT id, name FROM genres ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type genre struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var genres []genre
	for rows.Next() {
		var g genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		genres = append(genres, g)
	}
	if genres == nil {
		genres = []genre{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}
