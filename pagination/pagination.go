package pagination

import "gorm.io/gorm"

const globalDefaultPerPage = 30

// Pagination 用來表示分頁
type Pagination struct {
	Page       int `json:"page"`       // 目前頁面
	PerPage    int `json:"perPage"`    // 每頁顯示多少筆
	TotalCount int `json:"totalCount"` // 總筆數
	TotalPage  int `json:"totalPage"`  // 總頁數
}

func NewPagination(page int, perPage int) *Pagination {
	return &Pagination{
		Page:    page,
		PerPage: perPage,
	}
}

// LimitAndOffset return gorm scope function
func (p *Pagination) LimitAndOffset(db *gorm.DB) *gorm.DB {
	limit := p.PerPage
	offset := p.offset()

	return db.Limit(limit).Offset(offset)
}

// SetTotalCountAndPage 用來計算總數和分頁
func (p *Pagination) SetTotalCountAndPage(total int) {
	p.CheckOrSetDefault()
	p.TotalCount = total

	quotient := p.TotalCount / p.PerPage
	remainder := p.TotalCount % p.PerPage
	if remainder > 0 {
		quotient++
	}
	p.TotalPage = quotient
}

// CheckOrSetDefault 檢查PerPage值若未設置則設置預設值
func (p *Pagination) CheckOrSetDefault(params ...int) *Pagination {
	var defaultPerPage int
	if len(params) >= 1 {
		defaultPerPage = params[0]
	}

	if defaultPerPage <= 0 {
		defaultPerPage = globalDefaultPerPage
	}

	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = defaultPerPage
	}
	return p
}

// offset 計算 offset 的值
func (p *Pagination) offset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}
