package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/0501jp/shinjuku-lunch/api/db"
	"github.com/0501jp/shinjuku-lunch/api/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	dbOk := db.Init()
	if dbOk {
		defer db.Close()
	} else {
		log.Println("Offline mode — serving from seed.json")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/list", http.StatusFound)
	})

	r.Route("/api", func(r chi.Router) {
		r.Get("/list", handler.ListPage)
		if dbOk {
			r.Get("/restaurants", handler.GetRestaurants)
			r.Get("/restaurants/{id}", handler.GetRestaurant)
			r.Get("/lunch-logs", handler.GetLunchLogs)
			r.Post("/lunch-logs", handler.PostLunchLog)
			r.Get("/areas", handler.GetAreas)
			r.Get("/genres", handler.GetGenres)
		} else {
			h := handler.NewDevHandler()
			r.Get("/restaurants", h.GetRestaurants)
			r.Get("/restaurants/{id}", h.GetRestaurant)
			r.Get("/lunch-logs", h.GetLunchLogs)
			r.Post("/lunch-logs", h.PostLunchLog)
			r.Get("/areas", h.GetAreas)
			r.Get("/genres", h.GetGenres)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
