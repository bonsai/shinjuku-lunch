package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var DB *sql.DB
var Connected bool
var IsPostgres bool
var IsSQLite bool

func Init() bool {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// priority chain: Neon (PG) → MSSQL → SQLite → JSON
		if tryDSN("postgres://postgres:postgres@localhost:5432/shinjuku_lunch?sslmode=disable", true) {
			return true
		}
		log.Println("Neon not available, trying local MSSQL...")
		if tryLocalMSSQL() {
			return true
		}
		log.Println("MSSQL not available, trying SQLite...")
		if initSQLite("") {
			return true
		}
		log.Println("SQLite not available — falling back to seed.json")
		return false
	}

	// explicit DATABASE_URL set
	dsnLower := strings.ToLower(dsn)
	switch {
	case strings.HasPrefix(dsnLower, "postgres://") || strings.HasPrefix(dsnLower, "postgresql://"):
		return initPG(dsn)
	case strings.HasPrefix(dsnLower, "sqlserver://"):
		return initMSSQL(dsn)
	case strings.HasPrefix(dsnLower, "sqlite://"):
		return initSQLite(strings.TrimPrefix(dsn, "sqlite://"))
	}

	// unrecognized: try as-is, detect by content
	log.Printf("Unrecognized DATABASE_URL, trying as-is...")
	return tryDSN(dsn, strings.Contains(dsn, "postgres") || strings.Contains(dsn, "5432"))
}

func tryDSN(dsn string, isPG bool) bool {
	if isPG {
		return initPG(dsn)
	}
	return initMSSQL(dsn)
}

func tryLocalMSSQL() bool {
	connStr := "Server=localhost;Database=shinjuku_lunch;Trusted_Connection=True;TrustServerCertificate=True;"
	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return false
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return false
	}
	DB = db
	Connected = true
	IsPostgres = false
	IsSQLite = false
	initMSSQLSchema()
	fmt.Println("Connected to SQL Server (local)")
	return true
}

// P returns the parameter placeholder for the current DB driver.
//   PG ($1..$N), MSSQL (@p1..@pN), SQLite (?)
func P(n int) string {
	switch {
	case IsPostgres:
		return fmt.Sprintf("$%d", n)
	case IsSQLite:
		return "?"
	default:
		return fmt.Sprintf("@p%d", n)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
