package condition

import (
	"fmt"

	"gorm.io/gorm"
)

type Sorting struct {
	SortField string
	SortOrder string
}

func (s *Sorting) Sort(db *gorm.DB) *gorm.DB {
	if len(s.SortField) != 0 && len(s.SortOrder) != 0 {
		db = db.Order(fmt.Sprintf("%s %s", s.SortField, s.SortOrder))
	}

	return db
}
