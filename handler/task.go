package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/webdevfuel/projectmotor/template"
	"github.com/webdevfuel/projectmotor/validator"
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

type CreateTaskForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	ProjectID   string `form:"project_id"`
}

func (data CreateTaskForm) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&data.ProjectID, validation.Required, is.Digit),
	)
}

func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var data CreateTaskForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	projects, err := h.ProjectService.GetAll(h.GetUserFromContext(r.Context()).ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.TaskNewForm(errors, projects)
		err = component.Render(r.Context(), w)
		if err != nil {
			fail(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	projectID, err := strconv.Atoi(data.ProjectID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	user := h.GetUserFromContext(r.Context())
	err = h.TaskService.Create(data.Title, data.Description, int32(projectID), user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Redirect", "http://localhost:3000/tasks")
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	tasks, err := h.TaskService.GetAll(user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.Tasks(tasks)
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	task, err := h.TaskService.Get(int32(id), user.ID)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
	component := template.TaskEditForm(task, validator.NewValidatedSlice())
	err = component.Render(r.Context(), w)
	if err != nil {
		fail(w, err, http.StatusInternalServerError)
		return
	}
}
