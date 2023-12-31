package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/db"
)

type Handler struct {
	userService *db.UserService
}

type HandlerOptions struct {
	DB *sqlx.DB
}

func NewHandler(options HandlerOptions) *Handler {
	userService := db.NewUserService(options.DB)
	return &Handler{
		userService: userService,
	}
}
