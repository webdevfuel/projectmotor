package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevfuel/projectmotor/test"
)

func TestTask(t *testing.T) {
	handler, server := test.NewServer()
	defer server.Close()

	cookie, err := test.SetUserSession(server, 1)
	if err != nil {
		t.Errorf("error setting user test session %s", err)
		return
	}

	err = test.ResetAndSeedDB(handler.DB)
	if err != nil {
		t.Errorf("error resetting and seeding database %s", err)
		return
	}

	t.Run("navigating to tasks page lists all tasks", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks")),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		body := test.Body(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// body assertions
		assert.Contains(body, "Task 1")
		assert.Contains(body, "Task 2")
		assert.Contains(body, "Task 3")
		assert.Contains(body, "Task 4")
		assert.NotContains(body, "Task 5")
		assert.NotContains(body, "Task 6")
		assert.NotContains(body, "Task 7")
		assert.NotContains(body, "Task 8")
	})

	t.Run("navigating to tasks page with id lists all tasks of project", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks?project=1")),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		body := test.Body(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// body assertions
		assert.Contains(body, "Task 1")
		assert.Contains(body, "Task 2")
		assert.NotContains(body, "Task 3")
		assert.NotContains(body, "Task 4")
	})

	t.Run("edit task displays form and data", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks/1/edit")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Get),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// form assertions
		form := test.NewForm(doc, "task-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("Task 1", title.Value)
		assert.Equal("Title (required)", title.Label)

		description := form.MustGetFieldByID("description")
		assert.Equal("", description.Value)
		assert.Equal("Description", description.Label)

		assert.Equal(1, test.FindByText(doc, "button", "Save task").Size())
		assert.Equal(1, test.FindByText(doc, "button", "Delete task").Size())
	})

	t.Run("update task returns toast", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks/1")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Patch),
			test.WithFormValues(
				test.FormValue{
					Key:   "title",
					Value: "Updated Task 1",
				},
				test.FormValue{
					Key:   "description",
					Value: "",
				},
			),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// body assertions
		assert.Equal("Task updated successfully", doc.Find("div[id='toast'] p").Text())

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from tasks where title = 'Updated Task 1' and id = 1")
		assert.Equal(1, count)
	})

	t.Run("get task returns row", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks/1")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Get),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// body assertions
		assert.Equal("Updated Task 1", doc.Find("div p").Text())
		assert.Equal("Edit", doc.Find("div button").Text())
	})

	t.Run("navigating to tasks new page shows form", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks/new")),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// form assertions
		form := test.NewForm(doc, "task-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("", title.Value)
		assert.Equal("Title (required)", title.Label)

		description := form.MustGetFieldByID("description")
		assert.Equal("", description.Value)
		assert.Equal("Description", description.Label)

		projects := []string{}
		project := form.MustGetSelectByID("project_id")
		for _, option := range project.Options {
			if option.Value == "" {
				if !option.Selected {
					assert.Fail("empty value should be selected by default")
				}
			} else {
				projects = append(projects, option.Label)
			}
		}
		assert.Equal([]string{"Project 2", "Project 1"}, projects)

		assert.Equal(1, test.FindByText(doc, "button", "New task").Size())
	})

	t.Run("new task returns form with errors", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Post),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// form assertions
		form := test.NewForm(doc, "task-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("Title (required)", title.Label)
		assert.Equal("", title.Value)
		assert.Equal("cannot be blank", title.Error)

		description := form.MustGetFieldByID("description")
		assert.Equal("Description", description.Label)
		assert.Equal("", description.Value)
		assert.Equal("", description.Error)

		project := form.MustGetSelectByID("project_id")
		assert.Equal("Project", project.Label)
		assert.Equal("", project.Error)
		for _, option := range project.Options {
			if option.Selected {
				assert.Equal("", option.Value)
			}
		}

		assert.Equal(1, test.FindByText(doc, "button", "New task").Size())
	})

	t.Run("new task without project redirects to '/tasks'", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Post),
			test.WithFormValues(
				test.FormValue{
					Key:   "title",
					Value: "Task 999",
				},
				test.FormValue{
					Key:   "description",
					Value: "",
				},
				test.FormValue{
					Key:   "project_id",
					Value: "",
				},
			),
		)
		res := test.Do(req)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// redirection assertions
		assert.Equal("http://localhost:3000/tasks", res.Header.Get("Hx-Redirect"))

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from tasks where title = 'Task 999' and project_id is null")
		assert.Equal(1, count)
	})

	t.Run("new task with project redirects to '/tasks'", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "tasks")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Post),
			test.WithFormValues(
				test.FormValue{
					Key:   "title",
					Value: "Task 123",
				},
				test.FormValue{
					Key:   "description",
					Value: "",
				},
				test.FormValue{
					Key:   "project_id",
					Value: "1",
				},
			),
		)
		res := test.Do(req)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// redirection assertions
		assert.Equal("http://localhost:3000/tasks", res.Header.Get("Hx-Redirect"))

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from tasks where title = 'Task 123' and project_id = 1")
		assert.Equal(1, count)
	})
}
