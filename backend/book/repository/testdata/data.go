package testdata

import (
	"app/book/entity"
	"time"
)

var BookForListTestData = []entity.BookForList{
	entity.BookForList{
		BookID:      "bkN1A123456789A123456789A1",
		Title:       "ハリー・ポッターと賢者の石",
		CoverID:     "mageN1A23456789A123456789A",
		LastRegDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		ReviewCount: 5,
		Rating:      4.5,
	},
	entity.BookForList{
		BookID:      "bkN2A123456789A123456789A1",
		Title:       "ハリー・ポッターと秘密の部屋",
		CoverID:     "mageN2A23456789A123456789A",
		LastRegDate: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		ReviewCount: 6,
		Rating:      5,
	},
	entity.BookForList{
		BookID:      "bkN2A123456789A123456789A1",
		Title:       "プログラミングの基礎",
		CoverID:     "mageN2A23456789A123456789A",
		LastRegDate: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
		ReviewCount: 3,
		Rating:      3,
	},
}
