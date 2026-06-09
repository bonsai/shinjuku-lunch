package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

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
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')))`,
		`CREATE TABLE IF NOT EXISTS lunch_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			restaurant_id INTEGER REFERENCES restaurants(id),
			menu TEXT NOT NULL,
			price INTEGER NOT NULL,
			rating INTEGER CHECK (rating BETWEEN 1 AND 5),
			comment TEXT,
			revisit INTEGER DEFAULT 0,
			visited_date TEXT DEFAULT (date('now')),
			created_at TEXT DEFAULT (datetime('now')))`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Printf("SQLite schema init: %v", err)
		}
	}
}

func fileDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return ""
}
