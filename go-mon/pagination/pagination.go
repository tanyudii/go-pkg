package pagination

import (
	"math"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaximumLimit = 100
)

type Pagination struct {
	Page      int32
	Limit     int32
	Total     int64
	TotalPage int32
}

func NewPagination(page int32, limit int32) *Pagination {
	p := &Pagination{Page: page, Limit: limit}
	p.Validate()
	return p
}

func NewPaginationWithMax(page int32) *Pagination {
	p := &Pagination{Page: page, Limit: MaximumLimit}
	p.Validate()
	return p
}

func (p *Pagination) Validate() *Pagination {
	if p.Page <= 0 {
		p.Page = DefaultPage
	}
	if p.Limit <= 0 || p.Limit > MaximumLimit {
		p.Limit = DefaultLimit
	}
	return p
}

func (p *Pagination) SetPagination(total ...int64) {
	if len(total) != 0 {
		p.Total = total[0]
	}
	if p.Total == 0 {
		return
	}
	p.TotalPage = int32(math.Ceil(float64(p.Total) / float64(p.Limit)))
}

func (p *Pagination) GetOffset() int32 {
	var offset int32
	if p.Page > 0 {
		offset = p.Limit * (p.Page - 1)
	}
	return offset
}

func (p *Pagination) HasNextPage() bool {
	return p.Page < p.TotalPage
}

func (p *Pagination) Incr() {
	p.Page++
}
