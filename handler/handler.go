package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/webdevfuel/projectmotor/database"
)

type Handler struct {
	userService *database.UserService
}

type HandlerOptions struct {
	DB *sqlx.DB
}

func NewHandler(options HandlerOptions) *Handler {
	userService := database.NewUserService(options.DB)
	return &Handler{
		userService: userService,
	}
}
