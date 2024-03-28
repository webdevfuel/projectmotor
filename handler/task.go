package handler

import (
	"net/http"

	"github.com/webdevfuel/projectmotor/template"
)

func (h Handler) NewTask(w http.ResponseWriter, r *http.Request) {
	component := template.TaskNew()
	err := component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}
