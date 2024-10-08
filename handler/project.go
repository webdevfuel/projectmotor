package handler

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/webdevfuel/projectmotor/template"
	"github.com/webdevfuel/projectmotor/template/toast"
	"github.com/webdevfuel/projectmotor/validator"
)

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := template.Projects(projects)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) NewProject(w http.ResponseWriter, r *http.Request) {
	component := template.ProjectNew()
	err := component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
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

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var data CreateProjectForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.ProjectNewForm(errors)
		err = component.Render(r.Context(), w)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	user := h.GetUserFromContext(r.Context())
	_, err = h.ProjectService.Create(data.Title, data.Description, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.Redirect(w, "http://localhost:3000/projects")
}

func (h *Handler) EditProject(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := h.GetIDFromRequest(r, "id")
	project, err := h.ProjectService.Get(id, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := template.ProjectEdit(project)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ToggleProjectPublished(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := h.GetIDFromRequest(r, "id")
	project, err := h.ProjectService.TogglePublished(id, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.TriggerEvent(w, fmt.Sprintf("toggle-project-status:%d", project.ID))
	component := template.ProjectStatusLabel(project.Published)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component = toast.Toast(toast.ToastOpts{
		Message: "Project updated successfully",
		Type:    "success",
		SwapOOB: true,
	})
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
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

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	var data UpdateProjectForm
	id, _ := h.GetIDFromRequest(r, "id")
	user := h.GetUserFromContext(r.Context())
	project, err := h.ProjectService.Get(id, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.ProjectEditForm(project, errors, template.NewProjectEditFormOpts())
		err = component.Render(r.Context(), w)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	project, err = h.ProjectService.Update(id, data.Title, data.Description, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := template.ProjectEditForm(project, validator.NewValidatedSlice(), template.ProjectEditFormOpts{
		SwapOOB: true,
	})
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component = toast.Toast(toast.ToastOpts{
		Message: "Project updated successfully",
		Type:    "success",
		SwapOOB: true,
	})
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := h.GetIDFromRequest(r, "id")
	err := h.ProjectService.Delete(id, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.Redirect(w, "http://localhost:3000/projects")
}
