package repositories

import (
	"github.com/jmoiron/sqlx"
)

type DatabaseRepositories struct {
	*UserQueries
	*CatQueries
	*CatMatchQueries
}

func New(db *sqlx.DB) *DatabaseRepositories {
	return &DatabaseRepositories{
		UserQueries:     &UserQueries{DB: db},
		CatQueries:      &CatQueries{DB: db},
		CatMatchQueries: &CatMatchQueries{DB: db},
	}
}
