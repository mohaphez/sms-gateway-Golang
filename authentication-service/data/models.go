package data

import (
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const dbTimeout = time.Second * 3

var db *pgxpool.Pool

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.

func New(dbpool *pgxpool.Pool) Models {
	db = dbpool
	createTables()
	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

func createTables() {
	createUserstable()
}
