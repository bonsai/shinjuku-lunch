package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Seed struct {
	Areas       []Area       `json:"areas"`
	Genres      []Genre      `json:"genres"`
	Restaurants []Restaurant `json:"restaurants"`
	LunchLogs   []LunchLog   `json:"lunch_logs"`
}

type Area struct {
	Name string `json:"name"`
}

type Genre struct {
	Name string `json:"name"`
}

type Restaurant struct {
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

type LunchLog struct {
	Restaurant  string `json:"restaurant"`
	Menu        string `json:"menu"`
	Price       int    `json:"price"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment,omitempty"`
	Revisit     bool   `json:"revisit"`
	VisitedDate string `json:"visited_date"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ctx := context.Background()

	seedPath := os.Getenv("SEED_JSON")
	if seedPath == "" {
		seedPath = "../neon/seed.json"
	}

	if data, err := os.ReadFile(seedPath); err == nil {
		importJSON(ctx, pool, data)
		return
	}
	log.Println("seed.json not found, trying raw/ md import")

	rawDir := os.Getenv("RAW_DIR")
	if rawDir == "" {
		rawDir = "../../DATA/raw"
	}
	importMarkdown(ctx, pool, rawDir)
}

func importJSON(ctx context.Context, pool *pgxpool.Pool, data []byte) {
	var seed Seed
	if err := json.Unmarshal(data, &seed); err != nil {
		log.Fatalf("JSON parse error: %v", err)
	}

	for _, a := range seed.Areas {
		pool.Exec(ctx, "INSERT INTO areas (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", a.Name)
	}
	fmt.Printf("areas: %d\n", len(seed.Areas))

	for _, g := range seed.Genres {
		pool.Exec(ctx, "INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", g.Name)
	}
	fmt.Printf("genres: %d\n", len(seed.Genres))

	for _, r := range seed.Restaurants {
		var areaID, genreID int
		pool.QueryRow(ctx, "SELECT id FROM areas WHERE name=$1", r.Area).Scan(&areaID)
		pool.QueryRow(ctx, "SELECT id FROM genres WHERE name=$1", r.Genre).Scan(&genreID)

		_, err := pool.Exec(ctx, `
			INSERT INTO restaurants (name, area_id, genre_id, address, station, walk_min, latitude, longitude, business_hours, url_tabelog, url_hotpepper, notes)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
			ON CONFLICT DO NOTHING`,
			r.Name, areaID, genreID, r.Address, r.Station, r.WalkMin,
			r.Latitude, r.Longitude, r.BusinessHours, r.URLTabelog, r.URLHotpepper, r.Notes)
		if err != nil {
			log.Printf("  skip restaurant %s: %v", r.Name, err)
		}
	}
	fmt.Printf("restaurants: %d\n", len(seed.Restaurants))

	for _, l := range seed.LunchLogs {
		var restID int
		err := pool.QueryRow(ctx, "SELECT id FROM restaurants WHERE name=$1", l.Restaurant).Scan(&restID)
		if err != nil {
			log.Printf("  skip log for %s: not found", l.Restaurant)
			continue
		}
		pool.Exec(ctx, `
			INSERT INTO lunch_logs (restaurant_id, menu, price, rating, comment, revisit, visited_date)
			VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			restID, l.Menu, l.Price, l.Rating, l.Comment, l.Revisit, l.VisitedDate)
	}
	fmt.Printf("lunch_logs: %d\n", len(seed.LunchLogs))
	fmt.Println("Seed import complete (JSON)")
}

func importMarkdown(ctx context.Context, pool *pgxpool.Pool, rawDir string) {
	patterns := []string{"*ランチ*", "*lunch*", "*めも*"}
	areaRe := regexp.MustCompile(`エリア\s*[|]\s*(.+?)\s*[|]`)
	priceRe := regexp.MustCompile(`[¥￥](\d{3,4})`)
	nameRe := regexp.MustCompile(`[|]\s*店名\s*[|]\s*(.+?)\s*[|]`)
	menuRe := regexp.MustCompile(`[|]\s*注文\s*[|]\s*(.+?)\s*[|]`)
	ratingRe := regexp.MustCompile(`[|]\s*評価\s*[|]\s*(.+?)\s*[|]`)
	revisitRe := regexp.MustCompile(`[|]\s*再訪\s*[|]\s*(.+?)\s*[|]`)
	commentRe := regexp.MustCompile(`[|]\s*一言\s*[|]\s*(.+?)\s*[|]`)

	count := 0
	err := filepath.WalkDir(rawDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		matched := false
		for _, p := range patterns {
			if matched, _ = filepath.Match(p, name); matched {
				break
			}
		}
		if !matched {
			return nil
		}

		fmt.Println("Processing:", path)
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		var lines []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		text := strings.Join(lines, "\n")

		matches := strings.Split(text, "##")
		for _, m := range matches {
			if parseEntry(ctx, pool, m, areaRe, priceRe, nameRe, menuRe, ratingRe, revisitRe, commentRe) {
				count++
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Imported %d entries from markdown\n", count)
}

func parseEntry(ctx context.Context, pool *pgxpool.Pool, text string, areaRe, priceRe, nameRe, menuRe, ratingRe, revisitRe, commentRe *regexp.Regexp) bool {
	nameMatch := nameRe.FindStringSubmatch(text)
	if nameMatch == nil {
		return false
	}
	name := strings.TrimSpace(nameMatch[1])
	if name == "" || name == "店名" {
		return false
	}

	area := "新宿"
	if m := areaRe.FindStringSubmatch(text); m != nil {
		area = strings.TrimSpace(m[1])
	}
	price := 0
	if m := priceRe.FindStringSubmatch(text); m != nil {
		price, _ = strconv.Atoi(m[1])
	}
	menu := ""
	if m := menuRe.FindStringSubmatch(text); m != nil {
		menu = strings.TrimSpace(m[1])
	}
	ratingStr := ""
	if m := ratingRe.FindStringSubmatch(text); m != nil {
		ratingStr = strings.TrimSpace(m[1])
	}
	rating := 0
	if strings.Contains(ratingStr, "★") {
		rating = strings.Count(ratingStr, "★")
	}
	revisitStr := ""
	if m := revisitRe.FindStringSubmatch(text); m != nil {
		revisitStr = strings.TrimSpace(m[1])
	}
	revisit := strings.Contains(revisitStr, "したい")
	comment := ""
	if m := commentRe.FindStringSubmatch(text); m != nil {
		comment = strings.TrimSpace(m[1])
	}
	if comment == "一言" {
		comment = ""
	}

	var areaID int
	err := pool.QueryRow(ctx, "INSERT INTO areas (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id", area).Scan(&areaID)
	if err != nil {
		log.Printf("  area error: %v", err)
		return false
	}

	genre := "その他"
	genreID := 1
	pool.QueryRow(ctx, "INSERT INTO genres (name) VALUES ($1) ON CONFLICT (name) DO NOTHING RETURNING id", genre).Scan(&genreID)

	var restID int
	err = pool.QueryRow(ctx,
		`INSERT INTO restaurants (name, area_id, genre_id, notes)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT DO NOTHING
		 RETURNING id`,
		name, areaID, genreID, comment,
	).Scan(&restID)
	if err != nil {
		return false
	}

	if price > 0 && menu != "" {
		pool.Exec(ctx,
			`INSERT INTO lunch_logs (restaurant_id, menu, price, rating, comment, revisit)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			restID, menu, price, rating, comment, revisit,
		)
	}

	fmt.Printf("  imported: %s (%s, ¥%d)\n", name, area, price)
	return true
}
