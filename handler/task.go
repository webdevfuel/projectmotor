package handler

import (
	"net/http"

	"github.com/webdevfuel/projectmotor/template"
)

func (h Handler) NewTask(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.TaskNew(projects)
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}
