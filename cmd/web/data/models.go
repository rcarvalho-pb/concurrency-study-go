package data

import (
	"database/sql"
	"time"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

type Models struct {
	User User
	Plan Plan
}

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
		Plan: Plan{},
	}
}
