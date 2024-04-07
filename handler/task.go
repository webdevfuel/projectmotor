package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/webdevfuel/projectmotor/database"
	"github.com/webdevfuel/projectmotor/template"
	"github.com/webdevfuel/projectmotor/template/toast"
	"github.com/webdevfuel/projectmotor/util"
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
	projectID, err := util.Atoi32(data.ProjectID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	user := h.GetUserFromContext(r.Context())
	err = h.TaskService.Create(data.Title, data.Description, projectID, user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	h.Redirect(w, "http://localhost:3000/tasks")
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	// get query param from url
	project := h.GetURLQuery(r, "project")
	// get user from context
	user := h.GetUserFromContext(r.Context())
	// initialize tasks slice
	var tasks []database.Task
	// initialize project id int
	var projectId int32
	// check if project is empty
	if !project.IsEmpty {
		id, err := util.Atoi32(project.Value)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		projectId = id
	}
	// get appropriate tasks based on project id
	if project.IsEmpty {
		t, err := h.TaskService.GetAll(user.ID)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		tasks = t
	} else {
		t, err := h.TaskService.GetAllByProjectID(user.ID, projectId)
		if err != nil {
			h.Error(w, err, http.StatusInternalServerError)
			return
		}
		tasks = t
	}
	// get all projects
	projects, err := h.ProjectService.GetAll(user.ID)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// filter based on project existing or not
	var filter string
	if !project.IsEmpty {
		filter = fmt.Sprintf("%d", projectId)
	}
	// check if is htmx request and
	// render component based on request being htmx or not
	var component templ.Component
	htmx := h.IsHTMXRequest(r)
	if htmx {
		component = template.TasksColumns(tasks, true)
	} else {
		component = template.Tasks(tasks, projects, filter)
	}
	if htmx && project.IsEmpty {
		h.ReplaceUrl(w, "/tasks")
	} else if htmx && !project.IsEmpty {
		h.ReplaceUrl(w, fmt.Sprintf("/tasks?project=%d", projectId))
	}
	// render component
	err = component.Render(r.Context(), w)
	if err != nil {
		h.Error(w, err, http.StatusInternalServerError)
		return
	}
	// override component with tasks filter and render it
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
	task, err := h.TaskService.Get(id, user.ID)
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
