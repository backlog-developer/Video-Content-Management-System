package models

import "time"

type Course struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CategoryID   int       `json:"category_id"`
	CreatedAt    time.Time `json:"created_at"`
	InstructorID int       `json:"instructor_id"`
	CreatedBy    int       `json:"created_by"`
}
