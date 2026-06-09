package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func initMSSQL(dsn string) bool {
	// sqlserver://user:pass@host?database=xxx&...
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		log.Printf("Cannot open MSSQL: %v — falling back to seed.json", err)
		return false
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		log.Printf("Cannot ping MSSQL: %v — falling back to seed.json", err)
		return false
	}
	DB = db
	Connected = true
	IsPostgres = false
	initMSSQLSchema()
	fmt.Println("Connected to SQL Server")
	return true
}

func initMSSQLSchema() {
	queries := []string{
		`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='areas' AND xtype='U')
		 CREATE TABLE areas (id INT IDENTITY(1,1) PRIMARY KEY, name NVARCHAR(200) NOT NULL UNIQUE)`,
		`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='genres' AND xtype='U')
		 CREATE TABLE genres (id INT IDENTITY(1,1) PRIMARY KEY, name NVARCHAR(200) NOT NULL UNIQUE)`,
		`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='restaurants' AND xtype='U')
		 CREATE TABLE restaurants (
		   id INT IDENTITY(1,1) PRIMARY KEY, name NVARCHAR(500) NOT NULL,
		   area_id INT REFERENCES areas(id), genre_id INT REFERENCES genres(id),
		   address NVARCHAR(500), station NVARCHAR(200), walk_min INT,
		   latitude FLOAT(53), longitude FLOAT(53),
		   business_hours NVARCHAR(500), url_tabelog NVARCHAR(1000), url_hotpepper NVARCHAR(1000), notes NVARCHAR(2000),
		   created_at DATETIME2 DEFAULT GETUTCDATE(), updated_at DATETIME2 DEFAULT GETUTCDATE())`,
		`IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='lunch_logs' AND xtype='U')
		 CREATE TABLE lunch_logs (
		   id INT IDENTITY(1,1) PRIMARY KEY, restaurant_id INT REFERENCES restaurants(id),
		   menu NVARCHAR(500) NOT NULL, price INT NOT NULL,
		   rating INT CHECK (rating BETWEEN 1 AND 5),
		   comment NVARCHAR(1000), revisit BIT DEFAULT 0,
		   visited_date DATE DEFAULT GETDATE(), created_at DATETIME2 DEFAULT GETUTCDATE())`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Printf("MSSQL schema init: %v", err)
		}
	}
}
