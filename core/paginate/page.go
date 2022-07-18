package paginate

import (
	"math"
)

type Page struct {
	Ascending bool
	Limit     int
	Page      int
	Max       int
	Count     int // Count is the number items.
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
		Max:       1,
	}
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Page) setMaxPage(maxPage int) {
	if maxPage <= 0 {
		p.Max = 1
	} else {
		p.Max = maxPage
	}
}

func (p *Page) SetMaxCount(count int) {
	p.Count = count
	p.setMaxPage(int(math.Ceil(float64(count) / float64(p.Limit))))
}
