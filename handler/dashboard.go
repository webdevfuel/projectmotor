package handler

import (
	"fmt"
	"net/http"

	"github.com/webdevfuel/projectmotor/template"
)

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	component := template.Dashboard(fmt.Sprintf("Welcome back, %s!", user.Email))
	component.Render(r.Context(), w)
}
