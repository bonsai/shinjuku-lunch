package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/0501jp/shinjuku-lunch/api/db"
	"github.com/go-chi/chi/v5"
)

func setupTestDB(t *testing.T) func() {
	t.Helper()
	f, err := os.CreateTemp("", "shinjuku-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	path := f.Name()

	os.Setenv("DATABASE_URL", "sqlite://"+path)
	if !db.Init() {
		t.Fatal("failed to init SQLite test DB")
	}

	// remove leftover data (in case of schema reuse)
	db.DB.Exec("DELETE FROM lunch_logs")
	db.DB.Exec("DELETE FROM restaurants")
	db.DB.Exec("DELETE FROM genres")
	db.DB.Exec("DELETE FROM areas")

	seed := []string{
		`INSERT INTO areas (id, name) VALUES (1, '歌舞伎町'), (2, '大久保')`,
		`INSERT INTO genres (id, name) VALUES (1, 'タイ料理'), (2, '韓国料理')`,
		`INSERT INTO restaurants (id, name, area_id, genre_id, walk_min, latitude, longitude, notes) VALUES
			(1, 'バンタイ', 1, 1, 2, 35.6958, 139.7012, '平日ランチ¥950'),
			(2, 'TOP トッポギ', 1, 2, 3, 35.6955, 139.7025, '韓国料理¥700'),
			(3, 'うま煮や', 2, 2, 0, 35.6990, 139.6980, '西京焼き定食¥700')`,
		`INSERT INTO lunch_logs (id, restaurant_id, menu, price, rating, comment, revisit) VALUES
			(1, 1, 'グリーンカレー', 950, 4, '美味しかった', 1),
			(2, 1, 'パッタイ', 850, 3, '普通', 0),
			(3, 2, 'トッポギ', 700, 5, '最高', 1)`,
	}
	for _, s := range seed {
		if _, err := db.DB.Exec(s); err != nil {
			t.Fatalf("seed: %v", err)
		}
	}

	return func() {
		db.DB.Close()
		db.DB = nil
		db.Connected = false
		db.IsSQLite = false
		os.Remove(path)
		os.Unsetenv("DATABASE_URL")
	}
}

func newTestRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Get("/restaurants", GetRestaurants)
		r.Get("/restaurants/{id}", GetRestaurant)
		r.Get("/lunch-logs", GetLunchLogs)
		r.Post("/lunch-logs", PostLunchLog)
		r.Get("/areas", GetAreas)
		r.Get("/genres", GetGenres)
	})
	return r
}

func bodyString(t *testing.T, resp *httptest.ResponseRecorder) string {
	t.Helper()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestGetRestaurants(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/restaurants", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 3 {
		t.Fatalf("expected 3 restaurants, got %d", len(body))
	}
}

func TestGetRestaurantsFilterByArea(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/restaurants?area=歌舞伎町", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body) != 2 {
		t.Fatalf("expected 2 restaurants in 歌舞伎町, got %d", len(body))
	}
}

func TestGetRestaurantsFilterByGenre(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/restaurants?genre=タイ料理", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body) != 1 {
		t.Fatalf("expected 1 タイ料理, got %d", len(body))
	}
}

func TestGetRestaurantDetail(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/restaurants/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if body["name"] != "バンタイ" {
		t.Fatalf("expected バンタイ, got %v", body["name"])
	}
	logs, ok := body["logs"].([]interface{})
	if !ok {
		t.Fatal("expected logs array")
	}
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
}

func TestGetRestaurantNotFound(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/restaurants/999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d: %s", rec.Code, bodyString(t, rec))
	}
}

func TestPostLunchLog(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	body := `{"restaurant_id":1,"menu":"テストメニュー","price":500,"rating":4,"comment":"テスト","revisit":true}`
	req := httptest.NewRequest("POST", "/api/lunch-logs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp["menu"] != "テストメニュー" {
		t.Fatalf("expected テストメニュー, got %v", resp["menu"])
	}
	if resp["rating"] != float64(4) {
		t.Fatalf("expected rating 4, got %v", resp["rating"])
	}
}

func TestPostLunchLogValidation(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	tests := []struct {
		name string
		body string
		code int
	}{
		{"missing restaurant_id", `{"menu":"a","price":500,"rating":3}`, 400},
		{"missing menu", `{"restaurant_id":1,"price":500,"rating":3}`, 400},
		{"invalid rating", `{"restaurant_id":1,"menu":"a","price":500,"rating":6}`, 400},
		{"negative price", `{"restaurant_id":1,"menu":"a","price":-1,"rating":3}`, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/lunch-logs", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tt.code {
				t.Errorf("expected %d, got %d: %s", tt.code, rec.Code, bodyString(t, rec))
			}
		})
	}
}

func TestGetAreas(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/areas", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body) != 2 {
		t.Fatalf("expected 2 areas, got %d", len(body))
	}
}

func TestGetGenres(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	router := newTestRouter()
	req := httptest.NewRequest("GET", "/api/genres", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", rec.Code, bodyString(t, rec))
	}

	var body []map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &body)
	if len(body) != 2 {
		t.Fatalf("expected 2 genres, got %d", len(body))
	}
}
