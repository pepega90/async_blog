package models

import "time"

type Post struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}