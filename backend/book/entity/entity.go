package entity

import "time"

type Book struct {
	BookId        string
	IsBn          string
	Title         string
	Author        string
	Publisher     string
	PublishedDate string
	CoverId       string
}

type BookForList struct {
	BookId      string    `json:"bookId"`
	Title       string    `json:"title"`
	CoverId     string    `json:"coverId"`
	LastRegDate time.Time `json:"lastRegDate"`
	ReviewCount uint32    `json:"reviewCount"`
	Rating      float32   `json:"rating"`
}
