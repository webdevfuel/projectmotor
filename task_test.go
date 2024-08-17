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
}
