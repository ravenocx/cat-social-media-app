package cmd

import "github.com/jmoiron/sqlx"

type Http struct {
	DB *sqlx.DB
}

type iHttp interface {
	StartApp()
}

func New(http *Http) iHttp {
	return http
}
