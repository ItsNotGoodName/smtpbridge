package paginate

import (
	"net/url"
	"strconv"
)

type Paginate struct {
	PageMin int
	Page    int
	PageMax int
	URL     *url.URL
	Param   string
}

func New(URL url.URL, param string, pageMin, page, pageMax int) Paginate {
	URL.Query().Add(param, strconv.Itoa(page))
	return Paginate{
		PageMin: pageMin,
		Param:   param,
		Page:    page,
		PageMax: pageMax,
		URL:     &URL,
	}
}

func (p Paginate) link(page int) string {
	vals := p.URL.Query()
	vals.Set(p.Param, strconv.Itoa(page))
	p.URL.RawQuery = vals.Encode()
	return p.URL.String()
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
