package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

// A User is where we keep the information about a person, with the unique
// point of data being the email address.
//
// The email address is not set by the user when signing up, but rather
// extracted from the OAuth provider if a user doesn't yet exist with
// the email address sent by the OAuth provider.
//
// table: "users"
type User struct {
	ID                int32
	Name              pgtype.Text
	Email             string
	GitHubAccessToken string `db:"gh_access_token"`
	GitHubUserID      int32  `db:"gh_user_id"`
}

// A UserService is a connection to the database with methods
// for interacting with the "users" table.
type UserService struct {
	db *sqlx.DB
}

// NewUserService returns a pointer to UserService.
func NewUserService(db *sqlx.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// GetAccountByID returns a User, reports whether the user exists inside
// the database with the given id, and returns an error from the Get method.
//
// sql.ErrNoRows isn't treated as an error, but rather makes the function
// report the user doesn't exist.
func (us UserService) GetUserByID(ID int32) (User, bool, error) {
	var user User
	query := "select * from users where id = $1"
	err := us.db.Get(&user, query, ID)
	if err != sql.ErrNoRows {
		if err != nil {
			return User{}, false, err
		}
		return user, true, nil
	}
	return User{}, false, nil
}

// CreateUser returns a User and returns an error from the Get method.
//
// If successful, it inserts a new row into the "users" table with the given data.
func (us UserService) CreateUser(
	tx *sqlx.Tx,
	email string,
	ghAccessToken string,
	ghUserId int32,
) (User, error) {
	var user User
	query := "insert into users (email, gh_access_token, gh_user_id) values ($1, $2, $3) returning *"
	err := tx.Get(&user, query, email, ghAccessToken, ghUserId)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (us UserService) UpdateUser(
	tx *sqlx.Tx,
	email string,
	ghAccessToken string,
	ghUserId int32,
) (User, error) {
	var user User
	query := "update users set email = $1, gh_access_token = $2 where gh_user_id = $3 returning *;"
	err := tx.Get(&user, query, email, ghAccessToken, ghUserId)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (us UserService) GetSharedUsers(projectId int32) ([]User, error) {
	var users []User
	err := us.db.Select(
		&users,
		"select users.* from projects_users left join users on projects_users.user_id = users.id where projects_users.project_id = $1",
		projectId,
	)
	if err != nil {
		return []User{}, err
	}
	return users, nil
}

func (us UserService) GetUserByEmail(email string) (User, error) {
	var user User
	query := "select * from users where email = $1"
	err := us.db.Get(&user, query, email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (us UserService) MustGetUserByID(id int32) (User, error) {
	var user User
	query := "select * from users where id = $1"
	err := us.db.Get(&user, query, id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (us UserService) UserExists(email int32) (bool, error) {
	var count int
	query := "select count(*) from users where email = $1"
	err := us.db.Get(&count, query, email)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (us UserService) UserExistsByGitHubID(gitHubUserId int32) (bool, error) {
	var count int
	err := us.db.Get(&count, "select count(*) from users where gh_user_id = $1;", gitHubUserId)
	if err != nil {
		return false, err
	}
	return count != 0, err
}

func (us UserService) GetUserBySessionToken(sessionToken string) (User, error) {
	var user User
	err := us.db.Get(
		&user,
		"select users.* from users left join sessions on users.id = sessions.user_id where sessions.token = $1;",
		sessionToken,
	)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
