package app

type Pagination struct {
	Page      int
	PageSize  int
	Total     int
	PageCount int
	HasPages  bool
	Pages     []int
	HasPrev   bool
	PagePrev  int
	HasNext   bool
	PageNext  int
}

func NewPagination(page int, size int, total int) *Pagination {
	return &Pagination{Page: page, PageSize: size, Total: total}
}

func (p *Pagination) Prepare() {
	p.PageCount = p.Total/p.PageSize + 1
	if p.PageCount > 1 {
		p.HasPages = true
	}
	if p.Page > 1 {
		p.HasPrev = true
		p.PagePrev = p.Page - 1
	}
	if p.Page < p.PageCount {
		p.HasNext = true
		p.PageNext = p.Page + 1
	}

	p.Pages = []int{}
	for i := 0; i < p.PageCount; i++ {
		p.Pages = append(p.Pages, i+1)
	}
}

func (p *Pagination) IsActive(page int) bool {
	if p.Page == page {
		return true
	}
	return false
}
