package handler

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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
	component := template.ProjectEditForm(
		project,
		validator.NewValidatedSlice(),
		template.ProjectEditFormOpts{
			SwapOOB: true,
		},
	)
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

func (h *Handler) ShareProject(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := h.GetIDFromRequest(r, "id")
	project, err := h.ProjectService.Get(id, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	users, err := h.UserService.GetSharedUsers(id)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	var u []template.ProjectShareUser
	for _, user := range users {
		u = append(u, template.ProjectShareUser{
			ID:    user.ID,
			Email: user.Email,
		})
	}
	component := template.ProjectShare(project, u)
	component.Render(r.Context(), w)
}

type ShareProjectByEmailForm struct {
	Email  string `form:"email"`
	Notify bool   `form:"notify"`
}

func (data ShareProjectByEmailForm) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Email, validation.Required, is.Email),
		validation.Field(&data.Notify),
	)
}

func (h *Handler) ShareProjectByEmail(w http.ResponseWriter, r *http.Request) error {
	var data ShareProjectByEmailForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		return h.RenderComponents(w, r, http.StatusInternalServerError)
	}
	projectId, err := h.GetIDFromRequest(r, "id")
	if err != nil {
		return h.RenderComponents(w, r, http.StatusInternalServerError)
	}
	projectShareFormComponent := template.ProjectShareForm(projectId, errors)
	if !ok {
		return h.RenderComponents(w, r, http.StatusBadRequest, projectShareFormComponent)
	}
	owner := h.GetUserFromContext(r.Context())
	project, err := h.ProjectService.Get(projectId, owner.ID)
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			projectShareFormComponent,
			defaultErrorToastComponent(),
		)
	}
	user, err := h.UserService.GetUserByEmail(data.Email)
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			projectShareFormComponent,
			errorToastComponent("We couldn't find a user with the email address you provided."),
		)
	}
	if owner.ID == user.ID {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			projectShareFormComponent,
			errorToastComponent("It's not possible to share a project with yourself."),
		)
	}
	exists, err := h.ProjectService.Share(project.ID, user.ID)
	if exists {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			projectShareFormComponent,
			errorToastComponent(
				"The user with the email address you provided already has access to the project.",
			),
		)

	}
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			projectShareFormComponent,
			defaultErrorToastComponent(),
		)
	}
	return h.RenderComponents(
		w,
		r,
		http.StatusCreated,
		projectShareFormComponent,
		template.ProjectCurrentlySharedRow(projectId, template.ProjectShareUser{
			ID:    user.ID,
			Email: user.Email,
		}, true),
		successToastComponent("Project shared successfully."),
	)
}

func (h *Handler) RevokeProjectById(w http.ResponseWriter, r *http.Request) error {
	projectId, err := h.GetIDFromRequest(r, "projectId")
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			defaultErrorToastComponent(),
		)
	}
	userId, err := h.GetIDFromRequest(r, "userId")
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			defaultErrorToastComponent(),
		)
	}
	owner := h.GetUserFromContext(r.Context())
	project, err := h.ProjectService.Get(projectId, owner.ID)
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			defaultErrorToastComponent(),
		)
	}
	user, err := h.UserService.MustGetUserByID(userId)
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			defaultErrorToastComponent(),
		)
	}
	err = h.ProjectService.Revoke(project.ID, user.ID)
	if err != nil {
		return h.RenderComponents(
			w,
			r,
			http.StatusInternalServerError,
			defaultErrorToastComponent(),
		)
	}
	return h.RenderComponents(
		w,
		r,
		http.StatusOK,
		successToastComponent("Project unshared successfully."),
	)
}
