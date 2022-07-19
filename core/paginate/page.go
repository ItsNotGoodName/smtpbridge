package paginate

import (
	"math"
)

type Page struct {
	Ascending bool
	Limit     int
	Prev      int
	Page      int
	Next      int
	Min       int
	Max       int
	Count     int // Count is the number items.
}

func NewPage(page, limit int, ascending bool) Page {
	if limit <= 0 || limit > 100 {
		limit = 5
	}
	prev := page - 1
	if page < 1 {
		page = 1
		prev = 0
	}
	return Page{
		Ascending: ascending,
		Limit:     limit,
		Prev:      prev,
		Page:      page,
		Next:      0,
		Min:       1,
		Max:       0,
		Count:     0,
	}
}

func (p *Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Page) setMax(max int) {
	if max < 1 {
		p.Max = 1
		return
	}

	p.Max = max
	if p.Page < p.Max {
		p.Next = p.Page + 1
	}
}

func (p *Page) SetCount(count int) {
	p.Count = count
	p.setMax(int(math.Ceil(float64(count) / float64(p.Limit))))
}
