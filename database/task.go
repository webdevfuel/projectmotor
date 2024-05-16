package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

// A Task is a way for users to keep a title and helpful description of a
// thing they want to do, along with a completed state.
//
// table: "tasks"
type Task struct {
	ID          int32            `db:"id"`
	Title       string           `db:"title"`
	Description pgtype.Text      `db:"description"`
	OwnerID     int32            `db:"owner_id"`
	ProjectID   int32            `db:"project_id"`
	CreatedAt   pgtype.Timestamp `db:"created_at"`
	UpdatedAt   pgtype.Timestamp `db:"updated_at"`
}

// A TaskService is a connection to the database with methods
// for interacting with the "tasks" table.
type TaskService struct {
	db *sqlx.DB
}

// NewTaskService returns a pointer to TaskService.
func NewTaskService(db *sqlx.DB) *TaskService {
	return &TaskService{
		db: db,
	}
}

// Create returns a Task and returns an error from the Get method.
//
// If successful, it inserts a new row into the "tasks" table with the given data.
func (s *TaskService) Create(title string, description string, projectID int32, ownerID int32) error {
	var task Task
	return s.db.Get(&task, `
		INSERT INTO tasks (title, description, owner_id, project_id)
		    VALUES ($1, $2, $3, $4)
		RETURNING
		    *
	`, title, description, ownerID, projectID)
}

// GetAll returns a slice of Task and returns an error from the Select method.
func (s *TaskService) GetAll(ownerID int32) ([]Task, error) {
	var tasks []Task
	err := s.db.Select(&tasks, `
		SELECT
		    *
		FROM
		    tasks
		WHERE
		    owner_id = $1
	`, ownerID)
	return tasks, err
}

// GetAll returns a slice of Task and returns an error from the Select method.
//
// It filters the query by the given project id.
func (s *TaskService) GetAllByProjectID(ownerID int32, projectID int32) ([]Task, error) {
	var tasks []Task
	err := s.db.Select(&tasks, `
		SELECT
		    *
		FROM
		    tasks
		WHERE
		    owner_id = $1
		    AND project_id = $2
	`, ownerID, projectID)
	return tasks, err
}

// Get returns a Task and returns an error from the Get method.
func (s *TaskService) Get(taskID int32, ownerID int32) (Task, error) {
	var task Task
	err := s.db.Get(&task, `
		SELECT
		    *
		FROM
		    tasks
		WHERE
		    id = $1
		    AND owner_id = $2
	`, taskID, ownerID)
	return task, err
}

// Update returns a Task and returns an error from the Get method.
//
// If successful, it updates the "tasks" table row that matches the
// given task id and owner id, with the given title and description.
func (s *TaskService) Update(taskID int32, ownerID int32, title string, description string) error {
	_, err := s.db.Exec(`
		UPDATE
		    tasks
		SET
		    title = $3,
		    description = $4
		WHERE
		    id = $1
		    AND owner_id = $2
	`, taskID, ownerID, title, description)
	return err
}
