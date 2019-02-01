package models

import "time"

type Users struct {
	UserID    int       `json:"userId"`
	LocalID   string    `json:"localId"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}
