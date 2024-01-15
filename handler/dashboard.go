package handler

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/webdevfuel/projectmotor/template"
)

func (h Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	template.Dashboard.ExecuteWriter(pongo2.Context{"message": fmt.Sprintf("Welcome back, %s!", user.Email)}, w)
}
