package pagination

import (
	"math"
)

type Page struct {
	page    int
	perPage int
}

func NewPage(page int, perPage int) Page {
	if page < 1 {
		page = 1
	}
	if perPage < 5 {
		perPage = 5
	}
	if perPage > 100 {
		perPage = 100
	}

	return Page{page: page, perPage: perPage}
}

func (p Page) Offset() int {
	return (p.page - 1) * p.perPage
}

func (p Page) Limit() int {
	return p.perPage
}

type PageResult struct {
	Page       int
	PerPage    int
	TotalPages int
	TotalItems int
}

func NewPageResult(p Page, totalItems int) PageResult {
	totalPage := int(math.Ceil(float64(float64(totalItems) / float64(p.perPage))))
	if totalPage == 0 {
		totalPage = 1
	}
	return PageResult{
		Page:       p.page,
		PerPage:    p.perPage,
		TotalPages: totalPage,
		TotalItems: totalItems,
	}
}

func (p PageResult) Overflow() bool {
	return p.Page > p.TotalPages
}

func (p PageResult) HasNext() bool {
	return p.Page < p.TotalPages
}

func (p PageResult) Next() int {
	return p.Page + 1
}

func (p PageResult) HasPrevious() bool {
	return p.Page > 1
}

func (p PageResult) Previous() int {
	return p.Page - 1
}

func (p PageResult) Seen() int {
	seen := p.Page * p.PerPage
	if seen > p.TotalItems {
		return p.TotalItems
	}
	return seen
}
