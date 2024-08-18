package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/webdevfuel/projectmotor/util"
)

func Int4FromString(s string) (pgtype.Int4, error) {
	if s == "" {
		return pgtype.Int4{
			Int32: 0,
			Valid: false,
		}, nil
	}

	value, err := util.Atoi32(s)
	if err != nil {
		return pgtype.Int4{
			Int32: 0,
			Valid: false,
		}, err
	}

	return pgtype.Int4{
		Int32: value,
		Valid: true,
	}, nil
}
