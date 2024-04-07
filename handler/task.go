package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/template"
	"github.com/webdevfuel/projectmotor/template/toast"
	"github.com/webdevfuel/projectmotor/validator"
)

func (h *Handler) NewTask(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := template.TaskNew(projects)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
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

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var data CreateTaskForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	projects, err := h.ProjectService.GetAll(h.GetUserFromContext(r.Context()).ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.TaskNewForm(errors, projects)
		err = component.Render(r.Context(), w)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	projectID, err := strconv.Atoi(data.ProjectID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	user := h.GetUserFromContext(r.Context())
	err = h.TaskService.Create(data.Title, data.Description, int32(projectID), user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.Redirect(w, "http://localhost:3000/tasks")
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("project")
	user := h.GetUserFromContext(r.Context())
	var htmx bool
	if r.Header.Get("Hx-Request") != "" {
		htmx = true
	}
	var tasks []database.Task
	var projectID int
	if project == "" {
		t, err := h.TaskService.GetAll(user.ID)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		tasks = t
	} else {
		id, err := strconv.Atoi(project)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		t, err := h.TaskService.GetAllByProjectID(user.ID, int32(id))
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		tasks = t
		projectID = id
	}
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	var filter string
	if project != "" {
		filter = fmt.Sprintf("%d", projectID)
	} else {
		filter = ""
	}
	var component templ.Component
	if htmx {
		component = template.TasksColumns(tasks, true)
		if project == "" {
			h.ReplaceUrl(w, "/tasks")
		} else {
			h.ReplaceUrl(w, fmt.Sprintf("/tasks?project=%d", projectID))
		}
	} else {
		component = template.Tasks(tasks, projects, filter)
	}
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component = template.TasksFilter(filter)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserFromContext(r.Context())
	id, _ := h.GetIDFromRequest(r, "id")
	task, err := h.TaskService.Get(int32(id), user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.TriggerEvent(w, "open-modal")
	component := template.TaskEditForm(task, validator.NewValidatedSlice())
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}

type UpdateTaskForm struct {
	Title       string `form:"title"`
	Description string `form:"description"`
}

func (data UpdateTaskForm) Validate() error {
	return validation.ValidateStruct(&data,
		validation.Field(&data.Title, validation.Required, validation.Length(1, 255)),
	)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var data UpdateTaskForm
	ok, errors, err := validator.Validate(&data, r)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	taskId, _ := h.GetIDFromRequest(r, "id")
	userId := h.GetUserFromContext(r.Context()).ID
	task, err := h.TaskService.Get(taskId, userId)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		component := template.TaskEditForm(task, errors)
		err = component.Render(r.Context(), w)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		return
	}
	err = h.TaskService.Update(taskId, userId, data.Title, data.Description)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := toast.Toast(toast.ToastOpts{
		Message: "Task updated successfully",
		Type:    "success",
		SwapOOB: true,
	})
	h.TriggerEvent(w, fmt.Sprintf("update-task-row:%d", task.ID))
	h.Reswap(w, "none")
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskId, _ := h.GetIDFromRequest(r, "id")
	userId := h.GetUserFromContext(r.Context()).ID
	task, err := h.TaskService.Get(taskId, userId)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	component := template.TaskRow(task)
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
}
