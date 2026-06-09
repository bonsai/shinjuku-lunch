package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/fsnotify/fsnotify"
)

type devSeed struct {
	Areas       []devArea       `json:"areas"`
	Genres      []devGenre      `json:"genres"`
	Restaurants []devRestaurant `json:"restaurants"`
	LunchLogs   []devLog        `json:"lunch_logs"`
}

type devArea struct{ Name string `json:"name"` }
type devGenre struct{ Name string `json:"name"` }
type devRestaurant struct {
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
	CreatedAt     string   `json:"created_at"`
}
type devLog struct {
	ID          int    `json:"id"`
	Restaurant  string `json:"restaurant"`
	Menu        string `json:"menu"`
	Price       int    `json:"price"`
	Rating      int    `json:"rating"`
	Comment     string `json:"comment,omitempty"`
	Revisit     bool   `json:"revisit"`
	VisitedDate string `json:"visited_date"`
	CreatedAt   string `json:"created_at"`
}

type DevHandler struct {
	seed     devSeed
	nextLog  int
}

func NewDevHandler() *DevHandler {
	h := &DevHandler{nextLog: 1}
	h.loadSeed()
	go h.watchSeed()
	return h
}

func (h *DevHandler) loadSeed() {
	data, err := os.ReadFile("../neon/seed.json")
	if err != nil {
		log.Printf("DEV_MODE: cannot read seed.json: %v", err)
		return
	}
	var s devSeed
	if err := json.Unmarshal(data, &s); err != nil {
		log.Printf("DEV_MODE: cannot parse seed.json: %v", err)
		return
	}
	h.seed = s
	log.Println("seed.json loaded")
}

// watchSeed uses fsnotify to reload seed.json on write events.
func (h *DevHandler) watchSeed() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("DEV_MODE: fsnotify not available, using poll fallback: %v", err)
		h.pollSeed()
		return
	}
	defer watcher.Close()

	if err := watcher.Add("../neon/seed.json"); err != nil {
		log.Printf("DEV_MODE: cannot watch seed.json, using poll fallback: %v", err)
		h.pollSeed()
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Println("seed.json changed — reloading")
				h.loadSeed()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("fsnotify error: %v", err)
		}
	}
}

// pollSeed is a fallback when fsnotify is unavailable.
func (h *DevHandler) pollSeed() {
	var lastMod time.Time
	for {
		fi, err := os.Stat("../neon/seed.json")
		if err == nil && fi.ModTime().After(lastMod) {
			lastMod = fi.ModTime()
			h.loadSeed()
		}
		time.Sleep(5 * time.Second)
	}
}

func (h *DevHandler) GetAreas(w http.ResponseWriter, r *http.Request) {
	type item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var out []item
	for i, a := range h.seed.Areas {
		out = append(out, item{ID: i + 1, Name: a.Name})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DevHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	type item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var out []item
	for i, g := range h.seed.Genres {
		out = append(out, item{ID: i + 1, Name: g.Name})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DevHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)
	type restaurant struct {
		ID            int      `json:"id"`
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
		CreatedAt     string   `json:"created_at"`
	}
	var out []restaurant
	for i, r := range h.seed.Restaurants {
		out = append(out, restaurant{
			ID: i + 1, Name: r.Name, Area: r.Area, Genre: r.Genre,
			Address: r.Address, Station: r.Station, WalkMin: r.WalkMin,
			Latitude: r.Latitude, Longitude: r.Longitude,
			BusinessHours: r.BusinessHours, URLTabelog: r.URLTabelog,
			URLHotpepper: r.URLHotpepper, Notes: r.Notes, CreatedAt: now,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DevHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 || id > len(h.seed.Restaurants) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	rest := h.seed.Restaurants[id-1]
	type logItem struct {
		ID           int    `json:"id"`
		RestaurantID int    `json:"restaurant_id"`
		Menu         string `json:"menu"`
		Price        int    `json:"price"`
		Rating       int    `json:"rating"`
		Comment      string `json:"comment"`
		Revisit      bool   `json:"revisit"`
		VisitedDate  string `json:"visited_date"`
		CreatedAt    string `json:"created_at"`
	}
	var logs []logItem
	for _, l := range h.seed.LunchLogs {
		if l.Restaurant == rest.Name {
			logs = append(logs, logItem{
				ID: l.ID, RestaurantID: id, Menu: l.Menu,
				Price: l.Price, Rating: l.Rating, Comment: l.Comment,
				Revisit: l.Revisit, VisitedDate: l.VisitedDate, CreatedAt: l.VisitedDate,
			})
		}
	}
	now := time.Now().Format(time.RFC3339)
	resp := map[string]interface{}{
		"id": id, "name": rest.Name, "area": rest.Area, "genre": rest.Genre,
		"address": rest.Address, "station": rest.Station, "walk_min": rest.WalkMin,
		"latitude": rest.Latitude, "longitude": rest.Longitude,
		"business_hours": rest.BusinessHours, "url_tabelog": rest.URLTabelog,
		"url_hotpepper": rest.URLHotpepper, "notes": rest.Notes, "created_at": now,
		"logs": logs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *DevHandler) GetLunchLogs(w http.ResponseWriter, r *http.Request) {
	type item struct {
		ID           int    `json:"id"`
		RestaurantID int    `json:"restaurant_id"`
		Menu         string `json:"menu"`
		Price        int    `json:"price"`
		Rating       int    `json:"rating"`
		Comment      string `json:"comment"`
		Revisit      bool   `json:"revisit"`
		VisitedDate  string `json:"visited_date"`
		CreatedAt    string `json:"created_at"`
	}
	restaurantIDStr := r.URL.Query().Get("restaurant_id")
	var rid int
	if restaurantIDStr != "" {
		rid, _ = strconv.Atoi(restaurantIDStr)
	}
	var out []item
	for _, l := range h.seed.LunchLogs {
		if rid > 0 && (rid < 1 || rid > len(h.seed.Restaurants) || h.seed.Restaurants[rid-1].Name != l.Restaurant) {
			continue
		}
		restID := 0
		for j, r := range h.seed.Restaurants {
			if r.Name == l.Restaurant {
				restID = j + 1
				break
			}
		}
		out = append(out, item{
			ID: l.ID, RestaurantID: restID, Menu: l.Menu,
			Price: l.Price, Rating: l.Rating, Comment: l.Comment,
			Revisit: l.Revisit, VisitedDate: l.VisitedDate, CreatedAt: l.VisitedDate,
		})
	}
	if out == nil {
		out = []item{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *DevHandler) PostLunchLog(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RestaurantID int     `json:"restaurant_id"`
		Menu         string  `json:"menu"`
		Price        int     `json:"price"`
		Rating       int     `json:"rating"`
		Comment      *string `json:"comment,omitempty"`
		Revisit      *bool   `json:"revisit,omitempty"`
		VisitedDate  *string `json:"visited_date,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if input.RestaurantID < 1 || input.RestaurantID > len(h.seed.Restaurants) {
		http.Error(w, "restaurant_id not found", http.StatusBadRequest)
		return
	}
	revisit := false
	if input.Revisit != nil {
		revisit = *input.Revisit
	}
	visited := time.Now().Format("2006-01-02")
	if input.VisitedDate != nil && *input.VisitedDate != "" {
		visited = *input.VisitedDate
	}
	comment := ""
	if input.Comment != nil {
		comment = *input.Comment
	}
	h.nextLog++
	entry := devLog{
		ID: h.nextLog, Restaurant: h.seed.Restaurants[input.RestaurantID-1].Name,
		Menu: input.Menu, Price: input.Price, Rating: input.Rating,
		Comment: comment, Revisit: revisit, VisitedDate: visited,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	h.seed.LunchLogs = append(h.seed.LunchLogs, entry)
	type respItem struct {
		ID           int    `json:"id"`
		RestaurantID int    `json:"restaurant_id"`
		Menu         string `json:"menu"`
		Price        int    `json:"price"`
		Rating       int    `json:"rating"`
		Comment      string `json:"comment"`
		Revisit      bool   `json:"revisit"`
		VisitedDate  string `json:"visited_date"`
		CreatedAt    string `json:"created_at"`
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(respItem{
		ID: entry.ID, RestaurantID: input.RestaurantID, Menu: entry.Menu,
		Price: entry.Price, Rating: entry.Rating, Comment: entry.Comment,
		Revisit: entry.Revisit, VisitedDate: entry.VisitedDate, CreatedAt: entry.CreatedAt,
	})
}


