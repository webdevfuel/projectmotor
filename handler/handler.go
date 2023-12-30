package handler

type Handler struct{}

type HandlerOptions struct{}

func NewHandler(options HandlerOptions) *Handler {
	return &Handler{}
}
