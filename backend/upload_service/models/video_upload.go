package models

import "time"

//type RegisterRequest struct {
//Username string `json:"username"`
//Email    string `json:"email"`
//Password string `json:"password"`
//  Role     string `json:"role"` // optional, defaults to "user"
//}

//type LoginRequest struct {
// Username string `json:"username"`
//  Password string `json:"password"`
//}

type VideoUpload struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Filename     string    `json:"filename"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Duration     string    `json:"duration"` // optional
	UploadedBy   int       `json:"uploaded_by"`
	CourseID     int       `json:"course_id"`
	UploadStatus string    `json:"upload_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
