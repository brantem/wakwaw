package testutil

import (
	"example.com/shared/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewDB() (*db.DB, sqlmock.Sqlmock) {
	v, mock, _ := sqlmock.New()
	return &db.DB{Master: sqlx.NewDb(v, "sqlmock"), Read: sqlx.NewDb(v, "sqlmock")}, mock
}
