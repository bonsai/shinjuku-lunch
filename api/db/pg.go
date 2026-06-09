package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func initPG(dsn string) bool {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Cannot open PG: %v — falling back to seed.json", err)
		return false
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := withTimeout(10)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		log.Printf("Cannot ping PG: %v — falling back to seed.json", err)
		return false
	}
	DB = db
	Connected = true
	IsPostgres = true
	initPGSchema()
	fmt.Println("Connected to PostgreSQL")
	return true
}

func initPGSchema() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS areas (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE IF NOT EXISTS genres (id SERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE IF NOT EXISTS restaurants (
			id SERIAL PRIMARY KEY, name TEXT NOT NULL,
			area_id INTEGER REFERENCES areas(id), genre_id INTEGER REFERENCES genres(id),
			address TEXT, station TEXT, walk_min INTEGER,
			latitude DOUBLE PRECISION, longitude DOUBLE PRECISION,
			business_hours TEXT, url_tabelog TEXT, url_hotpepper TEXT, notes TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(), updated_at TIMESTAMPTZ DEFAULT NOW())`,
		`CREATE TABLE IF NOT EXISTS lunch_logs (
			id SERIAL PRIMARY KEY, restaurant_id INTEGER REFERENCES restaurants(id),
			menu TEXT NOT NULL, price INTEGER NOT NULL,
			rating INTEGER CHECK (rating BETWEEN 1 AND 5),
			comment TEXT, revisit BOOLEAN DEFAULT FALSE,
			visited_date DATE DEFAULT CURRENT_DATE, created_at TIMESTAMPTZ DEFAULT NOW())`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Printf("PG schema init: %v", err)
		}
	}
}
