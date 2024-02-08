package handler

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/webdevfuel/projectmotor/template"
	"github.com/webdevfuel/projectmotor/validator"
)

func (h Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	err = template.Projects.ExecuteWriter(pongo2.Context{
		"Projects": projects,
	}, w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) NewProject(w http.ResponseWriter, r *http.Request) {
	err := template.ProjectNew.ExecuteWriter(pongo2.Context{}, w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

type CreateProjectForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func (data CreateProjectForm) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, validation.Required, validation.Length(1, 255)),
	)
}

func (h Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var data CreateProjectForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		err = template.ProjectNewForm.ExecuteWriter(pongo2.Context{
			"Errors": errors,
		}, w)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	w.Header().Set("HX-Redirect", "http://localhost:3000/projects")
}
