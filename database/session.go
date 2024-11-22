package database

import (
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
)

type Session struct {
	ID        int32            `db:"id"`
	UserID    int32            `db:"user_id"`
	Token     string           `db:"token"`
	UserAgent string           `db:"user_agent"`
	CreatedAt pgtype.Timestamp `db:"created_at"`
}

type SessionService struct {
	db *sqlx.DB
}

func NewSessionService(db *sqlx.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

func (ss SessionService) GetSessionByToken(token string) (Session, bool, error) {
	var session Session
	query := "select * from sessions where token = $1"
	err := ss.db.Get(&session, query, token)
	if err != sql.ErrNoRows {
		if err != nil {
			return Session{}, false, err
		}
		return session, true, nil
	}
	return Session{}, false, nil
}

func (ss SessionService) CreateToken(
	tx *sqlx.Tx,
	userId int32,
	token string,
	userAgent string,
) error {
	_, err := tx.Exec(
		"insert into sessions (user_id, token, user_agent) values ($1, $2, $3);",
		userId,
		token,
		userAgent,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ss SessionService) DeleteToken(token string) error {
	_, err := ss.db.Exec(
		"delete from sessions where token = $1;",
		token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ss SessionService) DeleteAllTokens(userId int32, token string) error {
	_, err := ss.db.Exec(
		"delete from sessions where id in (select sessions.id from sessions left join users on users.id = sessions.user_id where users.id = $1 and sessions.token != $2);",
		userId,
		token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ss SessionService) GetAllSessions(userId int32) ([]Session, error) {
	var sessions []Session
	err := ss.db.Select(
		&sessions,
		"select sessions.* from sessions left join users on users.id = sessions.user_id where users.id = $1;",
		userId,
	)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
