package condition

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

var zeroTime = time.Time{}

type Where struct {
	SearchIn string  `json:"search_in"`
	Keyword  string  `json:"keyword"`
	IDs      []int64 `json:"ids"`

	CreatedAtLt  *time.Time `json:"create_at_lt"`
	CreatedAtLte *time.Time `json:"create_at_lte"`
	CreatedAtGt  *time.Time `json:"create_at_gt"`
	CreatedAtGte *time.Time `json:"create_at_gte"`

	UpdatedAtLt  *time.Time `json:"update_at_lt"`
	UpdatedAtLte *time.Time `json:"update_at_lte"`
	UpdatedAtGt  *time.Time `json:"update_at_gt"`
	UpdatedAtGte *time.Time `json:"update_at_gte"`

	DeletedAtLt  *time.Time `json:"delete_at_lt"`
	DeletedAtLte *time.Time `json:"delete_at_lte"`
	DeletedAtGt  *time.Time `json:"delete_at_gt"`
	DeletedAtGte *time.Time `json:"delete_at_gte"`
}

func (where *Where) Where(db *gorm.DB) *gorm.DB {
	if where.IDs != nil && len(where.IDs) != 0 {
		db = db.Where("id IN (?)", where.IDs)
	}

	if where.SearchIn != "" && where.Keyword != "" {
		fields := strings.Split(where.SearchIn, ",")
		for _, field := range fields {
			field := field
			db = db.Where(fmt.Sprintf("%s like ?", field), "%"+where.Keyword+"%")
		}
	}

	// time create
	if where.CreatedAtLt != nil && *where.CreatedAtLt != zeroTime {
		sec := (*where.CreatedAtLt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("created_at < ?", start)
	}

	if where.CreatedAtLte != nil && *where.CreatedAtLte != zeroTime {
		sec := (*where.CreatedAtLte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("created_at <= ?", start)
	}

	if where.CreatedAtGt != nil && *where.CreatedAtGt != zeroTime {
		sec := (*where.CreatedAtGt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("created_at > ?", start)
	}

	if where.CreatedAtGte != nil && *where.CreatedAtGte != zeroTime {
		sec := (*where.CreatedAtGte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("created_at >= ?", start)
	}

	// time update
	if where.UpdatedAtLt != nil && *where.UpdatedAtLt != zeroTime {
		sec := (*where.UpdatedAtLt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("updated_at < ?", start)
	}

	if where.UpdatedAtLte != nil && *where.UpdatedAtLte != zeroTime {
		sec := (*where.UpdatedAtLte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("updated_at <= ?", start)
	}

	if where.UpdatedAtGt != nil && *where.UpdatedAtGt != zeroTime {
		sec := (*where.UpdatedAtGt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("updated_at > ?", start)
	}

	if where.UpdatedAtGte != nil && *where.UpdatedAtGte != zeroTime {
		sec := (*where.UpdatedAtGte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("updated_at >= ?", start)
	}

	// time delete
	if where.DeletedAtLt != nil && *where.DeletedAtLt != zeroTime {
		sec := (*where.DeletedAtLt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("deleted_at < ?", start)
	}

	if where.DeletedAtLte != nil && *where.DeletedAtLte != zeroTime {
		sec := (*where.DeletedAtLte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("deleted_at <= ?", start)
	}

	if where.DeletedAtGt != nil && *where.DeletedAtGt != zeroTime {
		sec := (*where.DeletedAtGt).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("deleted_at > ?", start)
	}

	if where.DeletedAtGte != nil && *where.DeletedAtGte != zeroTime {
		sec := (*where.DeletedAtGte).Unix()
		start := time.Unix(int64(sec), 0).UTC()
		db = db.Where("deleted_at >= ?", start)
	}

	return db
}
