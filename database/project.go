package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

// A Project enables splitting tasks into different categories and
// sharing them all at once with other users.
//
// table: "projects"
type Project struct {
	// ID is the primary key of the "projects" table.
	ID int32 `db:"id"`
	// Title is the main identifier of a project.
	Title string `db:"title"`
	// Description is an optional text to describe the project.
	Description pgtype.Text `db:"description"`
	// Published reports whether the project is available with shared users.
	Published bool `db:"published"`
	// OwnerID is a foreign key to the "users" table.
	OwnerID int32 `db:"owner_id"`
	// CreatedAt tracks when the project was created by the owner.
	CreatedAt pgtype.Timestamp `db:"created_at"`
	// UpdatedAt tracks the last time the project was updated by the owner.
	UpdatedAt pgtype.Timestamp `db:"updated_at"`
	// Shared reports whether the project is shared or owned.
	// It should always be set by Go code, and doesn't map to
	// any column inside the "projects" table.
	Shared bool `db:"-"`
}

// An ProjectService is a connection to the database with methods
// for interacting with the "projects" table.
type ProjectService struct {
	db *sqlx.DB
}

// NewProjectService returns a pointer to ProjectService.
func NewProjectService(db *sqlx.DB) *ProjectService {
	return &ProjectService{
		db: db,
	}
}

// Create returns a Project and returns an error from the Get method.
//
// If successful, it inserts a new row into the "projects" table with
// the given title and description.
func (s ProjectService) Create(title string, description string, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "insert into projects (title, description, owner_id) values ($1, $2, $3) returning *", title, description, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// TogglePublished returns a Project and returns an error from the Get method.
//
// If successful, it updates the "published" column inside the "projects" table
// by the given id, to the opposite of the previous value.
func (s ProjectService) TogglePublished(projectID int32, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "update projects set published = not published where id = $1 and owner_id = $2 returning *", projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// Get returns a Project and returns an error from the Get method.
func (s ProjectService) Get(projectID int32, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "select * from projects where id = $1 and owner_id = $2", projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// GetAll returns a slice of Project and returns an error from the Select method.
func (s ProjectService) GetAll(ownerID int32) ([]Project, error) {
	var projects []Project
	query := "select * from projects where owner_id = $1 order by created_at desc"
	err := s.db.Select(&projects, query, ownerID)
	if err != nil {
		return []Project{}, err
	}
	return projects, nil
}

// Update returns a Project and returns an error from the Get method.
//
// If successful, it updates the "projects" table row that matches the
// given project id and owner id, with the given title and description.
func (s ProjectService) Update(projectID int32, title string, description string, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "update projects set title = $1, description = $2 where id = $3 and owner_id = $4 returning *", title, description, projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

// Delete returns an error from the Exec method.
//
// If successful, it delete the "projects" table row that matches the
// given project id and owner id.
func (s ProjectService) Delete(projectID int32, ownerID int32) error {
	_, err := s.db.Exec("delete from projects where id = $1 and owner_id = $2", projectID, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func (s ProjectService) Share(projectId int32, userId int32) (bool, error) {
	_, err := s.db.Exec("insert into projects_users (project_id, user_id) values ($1, $2);", projectId, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s ProjectService) Revoke(projectId int32, userId int32) error {
	_, err := s.db.Exec("delete from projects_users where project_id = $1 and user_id = $2;", projectId, userId)
	if err != nil {
		return err
	}
	return nil
}
