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

func (s ProjectService) Create(title string, description string, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "insert into projects (title, description, owner_id) values ($1, $2, $3) returning *", title, description, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (s ProjectService) TogglePublished(projectID int32, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "update projects set published = not published where id = $1 and owner_id = $2 returning *", projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (s ProjectService) Get(projectID int32, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "select * from projects where id = $1 and owner_id = $2", projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (s ProjectService) GetAll(ownerID int32) ([]Project, error) {
	var projects []Project
	query := "select * from projects where owner_id = $1 order by created_at desc"
	err := s.db.Select(&projects, query, ownerID)
	if err != nil {
		return []Project{}, err
	}
	return projects, nil
}

func (s ProjectService) Update(projectID int32, title string, description string, ownerID int32) (Project, error) {
	var project Project
	err := s.db.Get(&project, "update projects set title = $1, description = $2 where id = $3 and owner_id = $4 returning *", title, description, projectID, ownerID)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}
