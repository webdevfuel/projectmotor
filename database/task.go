package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID          int32            `db:"id"`
	Title       string           `db:"title"`
	Description pgtype.Text      `db:"description"`
	OwnerID     int32            `db:"owner_id"`
	ProjectID   int32            `db:"project_id"`
	CreatedAt   pgtype.Timestamp `db:"created_at"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at"`
}

type TaskService struct {
	db *sqlx.DB
}

func NewTaskService(db *sqlx.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

func (s *TaskService) Create(title string, description string, projectID int32, ownerID int32) error {
	var task Task
	return s.db.Get(&task, `
		INSERT INTO tasks (title, description, owner_id, project_id)
		    VALUES ($1, $2, $3, $4)
		RETURNING
		    *
	`, title, description, ownerID, projectID)
}