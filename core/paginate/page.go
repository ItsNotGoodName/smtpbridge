package paginate

import (
	"math"
)

type Page struct {
	Ascending bool
	Limit     int
	Page      int
	MaxPage   int // MaxPage is the maximum number of pages.
	MaxCount  int // MaxCount is the maximum number of items.
}

func NewPage(page, limit int, ascending bool) Page {
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	if page <= 0 {
		page = 1
	}
	return Page{
		Ascending: ascending,
		Limit:     limit,
		Page:      page,
		MaxPage:   1,
	}
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Page) SetMaxCount(count int) {
	p.MaxCount = count
	p.MaxPage = int(math.Ceil(float64(count) / float64(p.Limit)))
}
