package interfaces

import (
	"github.com/go-gorp/gorp"
	_ "github.com/mattn/go-sqlite3"
)

type AppContext struct {
	DB *gorp.DbMap
}

func NewAppContext() (*AppContext, error) {
	db, err := InitDb()
	if err != nil {
		return nil, err
	}

	return &AppContext{DB: db}, nil
}
