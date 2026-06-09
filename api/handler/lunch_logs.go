package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/0501jp/shinjuku-lunch/api/db"
	"github.com/0501jp/shinjuku-lunch/api/model"
)

func GetLunchLogs(w http.ResponseWriter, r *http.Request) {
	restaurantIDStr := r.URL.Query().Get("restaurant_id")

	query := "SELECT id, restaurant_id, menu, price, rating, comment, revisit, visited_date, created_at FROM lunch_logs WHERE 1=1"
	var args []interface{}
	n := 1

	if restaurantIDStr != "" {
		rid, err := strconv.Atoi(restaurantIDStr)
		if err == nil {
			query += " AND restaurant_id = " + db.P(n)
			args = append(args, rid)
			n++
		}
	}
	query += " ORDER BY visited_date DESC"

	rows, err := db.DB.QueryContext(r.Context(), query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []model.LunchLog
	for rows.Next() {
		var l model.LunchLog
		if err := rows.Scan(&l.ID, &l.RestaurantID, &l.Menu, &l.Price, &l.Rating,
			&l.Comment, &l.Revisit, &l.VisitedDate, &l.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []model.LunchLog{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func PostLunchLog(w http.ResponseWriter, r *http.Request) {
	var input model.LunchLogInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if input.RestaurantID == 0 || input.Menu == "" || input.Price < 0 || input.Rating < 1 || input.Rating > 5 {
		http.Error(w, "restaurant_id, menu, price (>=0), rating (1-5) are required", http.StatusBadRequest)
		return
	}

	revisit := false
	if input.Revisit != nil {
		revisit = *input.Revisit
	}
	visitedDate := time.Now().Format("2006-01-02")
	if input.VisitedDate != nil && *input.VisitedDate != "" {
		visitedDate = *input.VisitedDate
	}

	var l model.LunchLog
	var err error

	if db.IsPostgres {
		err = db.DB.QueryRowContext(r.Context(), `
			INSERT INTO lunch_logs (restaurant_id, menu, price, rating, comment, revisit, visited_date)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, restaurant_id, menu, price, rating, comment, revisit, visited_date, created_at`,
			input.RestaurantID, input.Menu, input.Price, input.Rating, input.Comment, revisit, visitedDate,
		).Scan(&l.ID, &l.RestaurantID, &l.Menu, &l.Price, &l.Rating,
			&l.Comment, &l.Revisit, &l.VisitedDate, &l.CreatedAt)
	} else if db.IsSQLite {
		// SQLite: RETURNING supported since 3.35
		err = db.DB.QueryRowContext(r.Context(), `
			INSERT INTO lunch_logs (restaurant_id, menu, price, rating, comment, revisit, visited_date)
			VALUES (?,?,?,?,?,?,?)
			RETURNING id, restaurant_id, menu, price, rating, comment, revisit, visited_date, created_at`,
			input.RestaurantID, input.Menu, input.Price, input.Rating, input.Comment, revisit, visitedDate,
		).Scan(&l.ID, &l.RestaurantID, &l.Menu, &l.Price, &l.Rating,
			&l.Comment, &l.Revisit, &l.VisitedDate, &l.CreatedAt)
	} else {
		err = db.DB.QueryRowContext(r.Context(), `
			INSERT INTO lunch_logs (restaurant_id, menu, price, rating, comment, revisit, visited_date)
			OUTPUT INSERTED.id, INSERTED.restaurant_id, INSERTED.menu, INSERTED.price,
			       INSERTED.rating, INSERTED.comment, INSERTED.revisit, INSERTED.visited_date, INSERTED.created_at
			VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7)`,
			input.RestaurantID, input.Menu, input.Price, input.Rating, input.Comment, revisit, visitedDate,
		).Scan(&l.ID, &l.RestaurantID, &l.Menu, &l.Price, &l.Rating,
			&l.Comment, &l.Revisit, &l.VisitedDate, &l.CreatedAt)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(l)
}
