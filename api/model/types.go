package model

import "time"

type Restaurant struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Area          string     `json:"area"`
	Genre         string     `json:"genre"`
	Address       *string    `json:"address,omitempty"`
	Station       *string    `json:"station,omitempty"`
	WalkMin       *int       `json:"walk_min,omitempty"`
	Latitude      *float64   `json:"latitude,omitempty"`
	Longitude     *float64   `json:"longitude,omitempty"`
	BusinessHours *string    `json:"business_hours,omitempty"`
	URLTabelog    *string    `json:"url_tabelog,omitempty"`
	URLHotpepper  *string    `json:"url_hotpepper,omitempty"`
	Notes         *string    `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type LunchLog struct {
	ID           int       `json:"id"`
	RestaurantID int       `json:"restaurant_id"`
	Menu         string    `json:"menu"`
	Price        int       `json:"price"`
	Rating       int       `json:"rating"`
	Comment      *string   `json:"comment,omitempty"`
	Revisit      bool      `json:"revisit"`
	VisitedDate  string    `json:"visited_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type LunchLogInput struct {
	RestaurantID int     `json:"restaurant_id"`
	Menu         string  `json:"menu"`
	Price        int     `json:"price"`
	Rating       int     `json:"rating"`
	Comment      *string `json:"comment,omitempty"`
	Revisit      *bool   `json:"revisit,omitempty"`
	VisitedDate  *string `json:"visited_date,omitempty"`
}

type Area struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
