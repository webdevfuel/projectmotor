package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type Project struct {
	ID          int32            `db:"id"`
	Title       string           `db:"title"`
	Description pgtype.Text      `db:"description"`
	Published   bool             `db:"published"`
	OwnerID     int32            `db:"owner_id"`
	CreatedAt   pgtype.Timestamp `db:"created_at"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at"`
	Shared      bool             `db:"-"`
}

type ProjectService struct {
	db *sqlx.DB
}

func NewProjectService(db *sqlx.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}
