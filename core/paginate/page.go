package paginate

import (
	"log"
	"math"
)

type Page struct {
	Ascending bool
	Limit     int
	Page      int
	MaxPages  int
}

func NewPage(page, limit int, ascending bool) Page {
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	if page <= 0 {
		page = 1
	}
	log.Println(page)
	return Page{
		Ascending: ascending,
		Limit:     limit,
		Page:      page,
		MaxPages:  1,
	}
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Page) SetCount(count int) {
	p.MaxPages = int(math.Ceil(float64(count) / float64(p.Limit)))
}
