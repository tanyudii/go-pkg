package mysql

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"pkg.tanyudii.me/go-pkg/go-mon/math"
	"pkg.tanyudii.me/go-pkg/go-mon/pagination"
	"pkg.tanyudii.me/go-pkg/go-mon/sort"
)

type Scope = func(db *gorm.DB) *gorm.DB

func Paginate(p *pagination.Pagination) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if p == nil {
			return db
		}
		return db.Offset(int(p.GetOffset())).Limit(int(p.Limit))
	}
}

func SearchLike(qb *gorm.DB, search string, columns []string) *gorm.DB {
	if search == "" || len(columns) == 0 {
		return qb
	}
	for i := range columns {
		qb = qb.Or(columns[i]+" LIKE ?", "%"+search+"%")
	}
	return qb
}

func SearchLikeRight(qb *gorm.DB, search string, columns []string) *gorm.DB {
	if search == "" || len(columns) == 0 {
		return qb
	}
	for i := range columns {
		qb = qb.Or(columns[i]+" LIKE ?", search+"%")
	}
	return qb
}

func SearchEqual(qb *gorm.DB, search string, columns []string) *gorm.DB {
	if search == "" || len(columns) == 0 {
		return qb
	}
	for i := range columns {
		qb = qb.Or(columns[i]+" = ?", search)
	}
	return qb
}

type MapSortableColumn map[int32]string

func Sort(sort int32, mapSortable MapSortableColumn) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if sort == 0 {
			return db
		}
		column, ok := mapSortable[math.AbsInt32(sort)]
		if !ok {
			return db
		}
		sortType := "ASC"
		if sort < 0 {
			sortType = "DESC"
		}
		return db.Order(fmt.Sprintf("%s %s", column, sortType))
	}
}

func SortBy(s *sort.Sort, mapSortable MapSortableColumn) Scope {
	return func(db *gorm.DB) *gorm.DB {
		if s == nil {
			return db
		}
		column, ok := mapSortable[s.By]
		if !ok {
			return db
		}
		sortType := "ASC"
		if s != nil && s.IsDesc() {
			sortType = "DESC"
		}
		return db.Order(fmt.Sprintf("%s %s", column, sortType))
	}
}

func Count(db *gorm.DB, model schema.Tabler) (total int64, err error) {
	err = db.Model(&model).Count(&total).Error
	return
}

func CountPg(db *gorm.DB, model schema.Tabler, pg *pagination.Pagination) (err error) {
	if pg == nil {
		return
	}
	pg.Total, err = Count(db, model)
	pg.SetPagination()
	return err
}

func Preload(relations []string) Scope {
	return func(db *gorm.DB) *gorm.DB {
		for i := range relations {
			db = db.Preload(relations[i])
		}
		return db
	}
}
