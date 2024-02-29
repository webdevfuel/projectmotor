package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	component := template.Projects(projects)
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) NewProject(w http.ResponseWriter, r *http.Request) {
	component := template.ProjectNew()
	err := component.Render(r.Context(), w)
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
		component := template.ProjectNewForm(errors)
		err = component.Render(r.Context(), w)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	user := h.GetUserFromContext(r.Context())
	_, err = h.ProjectService.Create(data.Title, data.Description, user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Redirect", "http://localhost:3000/projects")
}

func (h Handler) EditProject(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	project, err := h.ProjectService.Get(int32(id), user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.ProjectEdit(project)
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

func (h Handler) ToggleProjectPublished(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	project, err := h.ProjectService.TogglePublished(int32(id), user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.ProjectStatus(project)
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

type UpdateProjectForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func (data UpdateProjectForm) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, validation.Required, validation.Length(1, 255)),
	)
}

func (h Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var data UpdateProjectForm
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	user := h.GetUserFromContext(r.Context())
	project, err := h.ProjectService.Get(int32(id), user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.ProjectEditForm(project, errors, template.NewProjectEditFormOpts())
		err = component.Render(r.Context(), w)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	project, err = h.ProjectService.Update(int32(id), data.Title, data.Description, user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.ProjectEditForm(project, validator.NewValidatedSlice(), template.ProjectEditFormOpts{
		SwapOOB: true,
	})
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}
