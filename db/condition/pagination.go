package condition

import "github.com/jinzhu/gorm"

const globalDefaultPerPage = 30

// Pagination 用來表示分頁
type Pagination struct {
	Page       int64 `query:"page" json:"page" url:"page" description:"目前頁面"`
	PerPage    int64 `query:"perPage" json:"perPage" url:"perPage" description:"每頁顯示多少筆"`
	TotalCount int64 `json:"totalCount" url:"-" description:"總筆數"`
	TotalPage  int64 `json:"totalPage" url:"-" description:"總頁數"`
}

// Where return gorm scope function
func (p *Pagination) Where() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit, offset := p.LimitAndOffset()
		return db.Limit(limit).Offset(offset)
	}
}

// SetTotalCountAndPage 用來計算總數和分頁
func (p *Pagination) SetTotalCountAndPage(total int64) {
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
func (p *Pagination) CheckOrSetDefault(params ...int64) *Pagination {
	var defaultPerPage int64
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

// LimitAndOffset return limit and offset
func (p *Pagination) LimitAndOffset() (uint64, uint64) {
	return uint64(p.PerPage), uint64(p.Offset())
}

// Offset 計算 offset 的值
func (p *Pagination) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PerPage
}
