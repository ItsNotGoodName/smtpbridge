// Package paginate provides url links for a paginated resource.
package paginate

import (
	"net/url"
	"strconv"
)

type Paginate struct {
	PageMin int
	Page    int
	PageMax int
	Param   string
	url     *url.URL
}

func New(URL url.URL, param string, pageMin, page, pageMax int) Paginate {
	URL.Query().Add(param, strconv.Itoa(page))
	return Paginate{
		PageMin: pageMin,
		Param:   param,
		Page:    page,
		PageMax: pageMax,
		url:     &URL,
	}
}

func (p Paginate) link(page int) string {
	vals := p.url.Query()
	vals.Set(p.Param, strconv.Itoa(page))
	p.url.RawQuery = vals.Encode()
	return p.url.String()
}

func (p Paginate) HasPrev() bool {
	return p.Page > p.PageMin
}

func (p Paginate) HasNext() bool {
	return p.Page < p.PageMax
}

func (p Paginate) FirstLink() string {
	return p.link(p.PageMin)
}

func (p Paginate) LastLink() string {
	return p.link(p.PageMax)
}

func (p Paginate) PrevLink() string {
	return p.link(p.Page - 1)
}

func (p Paginate) NextLink() string {
	return p.link(p.Page + 1)
}
