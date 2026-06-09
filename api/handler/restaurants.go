package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0501jp/shinjuku-lunch/api/db"
	"github.com/0501jp/shinjuku-lunch/api/model"
	"github.com/go-chi/chi/v5"
)

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	area := r.URL.Query().Get("area")
	genre := r.URL.Query().Get("genre")
	priceMaxStr := r.URL.Query().Get("price_max")

	query := `SELECT r.id, r.name, a.name, g.name, r.address, r.station, r.walk_min,
		r.latitude, r.longitude, r.business_hours, r.url_tabelog,
		r.url_hotpepper, r.notes, r.created_at
	FROM restaurants r
	JOIN areas a ON r.area_id = a.id
	JOIN genres g ON r.genre_id = g.id WHERE 1=1`
	var args []interface{}
	n := 1

	if area != "" {
		query += fmt.Sprintf(" AND a.name = %s", db.P(n))
		args = append(args, area)
		n++
	}
	if genre != "" {
		query += fmt.Sprintf(" AND g.name = %s", db.P(n))
		args = append(args, genre)
		n++
	}
	if priceMaxStr != "" {
		priceMax, err := strconv.Atoi(priceMaxStr)
		if err == nil {
			query += fmt.Sprintf(" AND r.id IN (SELECT DISTINCT restaurant_id FROM lunch_logs WHERE price <= %s)", db.P(n))
			args = append(args, priceMax)
			n++
		}
	}
	query += " ORDER BY r.name"

	rows, err := db.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var restaurants []model.Restaurant
	for rows.Next() {
		var rest model.Restaurant
		if err := rows.Scan(
			&rest.ID, &rest.Name, &rest.Area, &rest.Genre,
			&rest.Address, &rest.Station, &rest.WalkMin,
			&rest.Latitude, &rest.Longitude, &rest.BusinessHours,
			&rest.URLTabelog, &rest.URLHotpepper, &rest.Notes, &rest.CreatedAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		restaurants = append(restaurants, rest)
	}
	if restaurants == nil {
		restaurants = []model.Restaurant{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`SELECT r.id, r.name, a.name, g.name, r.address, r.station, r.walk_min,
		r.latitude, r.longitude, r.business_hours, r.url_tabelog,
		r.url_hotpepper, r.notes, r.created_at
	FROM restaurants r
	JOIN areas a ON r.area_id = a.id
	JOIN genres g ON r.genre_id = g.id WHERE r.id = %s`, db.P(1))

	var rest model.Restaurant
	err = db.DB.QueryRowContext(r.Context(), query, id).Scan(
		&rest.ID, &rest.Name, &rest.Area, &rest.Genre,
		&rest.Address, &rest.Station, &rest.WalkMin,
		&rest.Latitude, &rest.Longitude, &rest.BusinessHours,
		&rest.URLTabelog, &rest.URLHotpepper, &rest.Notes, &rest.CreatedAt,
	)
	if err != nil {
		http.Error(w, "restaurant not found", http.StatusNotFound)
		return
	}

	logQuery := fmt.Sprintf("SELECT id, restaurant_id, menu, price, rating, comment, revisit, visited_date, created_at FROM lunch_logs WHERE restaurant_id = %s ORDER BY visited_date DESC", db.P(1))

	logRows, err := db.DB.QueryContext(r.Context(), logQuery, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer logRows.Close()

	var logs []model.LunchLog
	for logRows.Next() {
		var l model.LunchLog
		if err := logRows.Scan(&l.ID, &l.RestaurantID, &l.Menu, &l.Price, &l.Rating,
			&l.Comment, &l.Revisit, &l.VisitedDate, &l.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []model.LunchLog{}
	}

	resp := struct {
		model.Restaurant
		Logs []model.LunchLog `json:"logs"`
	}{rest, logs}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
