package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type seedArea struct{ Name string `json:"name"` }
type seedGenre struct{ Name string `json:"name"` }
type seedRestaurant struct {
	Name          string   `json:"name"`
	Area          string   `json:"area"`
	Genre         string   `json:"genre"`
	Address       *string  `json:"address,omitempty"`
	Station       *string  `json:"station,omitempty"`
	WalkMin       *int     `json:"walk_min,omitempty"`
	Latitude      *float64 `json:"latitude,omitempty"`
	Longitude     *float64 `json:"longitude,omitempty"`
	BusinessHours *string  `json:"business_hours,omitempty"`
	URLTabelog    *string  `json:"url_tabelog,omitempty"`
	URLHotpepper  *string  `json:"url_hotpepper,omitempty"`
	Notes         *string  `json:"notes,omitempty"`
}
type seedLog struct {
	Restaurant  string `json:"restaurant"`
	Menu        string `json:"menu"`
	Price       int    `json:"price"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment,omitempty"`
	Revisit     bool   `json:"revisit"`
	VisitedDate string `json:"visited_date"`
}
type seedData struct {
	Areas       []seedArea       `json:"areas"`
	Genres      []seedGenre      `json:"genres"`
	Restaurants []seedRestaurant `json:"restaurants"`
	LunchLogs   []seedLog        `json:"lunch_logs"`
}

const defaultSQLitePath = "./shinjuku_lunch.db"

func initSQLite(path string) bool {
	if path == "" {
		path = defaultSQLitePath
	}
	// ensure directory exists
	if dir := fileDir(path); dir != "" {
		os.MkdirAll(dir, 0755)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		log.Printf("Cannot open SQLite (%s): %v — falling back to seed.json", path, err)
		return false
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Printf("Cannot ping SQLite (%s): %v — falling back to seed.json", path, err)
		return false
	}

	// WAL mode for concurrent reads
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA busy_timeout=5000")
	db.Exec("PRAGMA foreign_keys=ON")

	DB = db
	Connected = true
	IsPostgres = false
	IsSQLite = true
	initSQLiteSchema()
	seedSQLiteFromJSON()
	fmt.Printf("Connected to SQLite (%s)\n", path)
	return true
}

func initSQLiteSchema() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS areas (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE IF NOT EXISTS restaurants (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			area_id INTEGER REFERENCES areas(id),
			genre_id INTEGER REFERENCES genres(id),
			address TEXT, station TEXT, walk_min INTEGER,
			latitude REAL, longitude REAL,
			business_hours TEXT, url_tabelog TEXT, url_hotpepper TEXT, notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP)`,
		`CREATE TABLE IF NOT EXISTS lunch_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			restaurant_id INTEGER REFERENCES restaurants(id),
			menu TEXT NOT NULL,
			price INTEGER NOT NULL,
			rating INTEGER CHECK (rating BETWEEN 1 AND 5),
			comment TEXT,
			revisit INTEGER DEFAULT 0,
			visited_date TEXT DEFAULT (date('now')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Printf("SQLite schema init: %v", err)
		}
	}
}

func seedSQLiteFromJSON() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM restaurants").Scan(&count)
	if count > 0 {
		return
	}

	data, err := os.ReadFile("../neon/seed.json")
	if err != nil {
		log.Printf("seed: cannot read ../neon/seed.json: %v", err)
		return
	}
	var s seedData
	if err := json.Unmarshal(data, &s); err != nil {
		log.Printf("seed: cannot parse seed.json: %v", err)
		return
	}

	areaMap := make(map[string]int64)
	for _, a := range s.Areas {
		res, err := DB.Exec("INSERT INTO areas (name) VALUES (?)", a.Name)
		if err != nil {
			log.Printf("seed: insert area %s: %v", a.Name, err)
			continue
		}
		id, _ := res.LastInsertId()
		areaMap[a.Name] = id
	}
	genreMap := make(map[string]int64)
	for _, g := range s.Genres {
		res, err := DB.Exec("INSERT INTO genres (name) VALUES (?)", g.Name)
		if err != nil {
			log.Printf("seed: insert genre %s: %v", g.Name, err)
			continue
		}
		id, _ := res.LastInsertId()
		genreMap[g.Name] = id
	}

	restNameToID := make(map[string]int64)
	for _, r := range s.Restaurants {
		aid, aok := areaMap[r.Area]
		gid, gok := genreMap[r.Genre]
		if !aok || !gok {
			log.Printf("seed: skip restaurant %s — unknown area/genre", r.Name)
			continue
		}
		now := time.Now().Format(time.RFC3339)
		res, err := DB.Exec(`INSERT INTO restaurants
			(name, area_id, genre_id, address, station, walk_min,
			 latitude, longitude, business_hours, url_tabelog, url_hotpepper, notes, created_at)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			r.Name, aid, gid, r.Address, r.Station, r.WalkMin,
			r.Latitude, r.Longitude, r.BusinessHours,
			r.URLTabelog, r.URLHotpepper, r.Notes, now)
		if err != nil {
			log.Printf("seed: insert restaurant %s: %v", r.Name, err)
			continue
		}
		id, _ := res.LastInsertId()
		restNameToID[r.Name] = id
	}

	for _, l := range s.LunchLogs {
		rid, ok := restNameToID[l.Restaurant]
		if !ok {
			log.Printf("seed: skip log for unknown restaurant %s", l.Restaurant)
			continue
		}
		revisit := 0
		if l.Revisit {
			revisit = 1
		}
		_, err := DB.Exec(`INSERT INTO lunch_logs
			(restaurant_id, menu, price, rating, comment, revisit, visited_date)
			VALUES (?,?,?,?,?,?,?)`,
			rid, l.Menu, l.Price, l.Rating, l.Comment, revisit, l.VisitedDate)
		if err != nil {
			log.Printf("seed: insert log for %s: %v", l.Restaurant, err)
		}
	}
	log.Printf("Seeded SQLite from seed.json (%d areas, %d genres, %d restaurants, %d logs)",
		len(s.Areas), len(s.Genres), len(s.Restaurants), len(s.LunchLogs))
}

func fileDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return ""
}
