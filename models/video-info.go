package models

type VideoInfo struct {
	ID       string `json:"id"`
	Title       string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Description string `json:"description"`
	// Add more fields as needed
}