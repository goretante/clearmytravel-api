package app

import (
	"github.com/goretante/clearmytravel-api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB      *pgxpool.Pool
	Queries *db.Queries
}
