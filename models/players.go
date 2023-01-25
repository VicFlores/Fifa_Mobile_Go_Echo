package models

import "time"

type Players struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	Club      string    `json:"club"`
	CreatedAt time.Time `json:"createdAt"`
}
