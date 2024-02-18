package entity

import "time"

type Book struct {
	BookID        string
	IsBn          string
	Title         string
	Author        string
	Publisher     string
	PublishedDate string
	CoverID       string
}

type BookForList struct {
	BookID      string    `json:"bookId"`
	Title       string    `json:"title"`
	CoverID     string    `json:"coverId"`
	LastRegDate time.Time `json:"lastRegDate"`
	ReviewCount uint32    `json:"reviewCount"`
	Rating      float32   `json:"rating"`
}
