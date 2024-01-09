package entity

import "time"

type Review struct {
	BookID   string
	ReviewID int
	UserID   string
	Rating   int8
	Review   string
	RegDate  time.Time
}
