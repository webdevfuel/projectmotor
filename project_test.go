package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webdevfuel/projectmotor/test"
)

func TestProject(t *testing.T) {
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

	t.Run("navigating to projects page shows lists all projects", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects")),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		body := test.Body(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// body assertions
		assert.Contains(body, "Project 1")
		assert.Contains(body, "Project 2")
		assert.NotContains(body, "Project 3")
		assert.NotContains(body, "Project 4")
	})

	t.Run("navigating to projects new page shows form", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects/new")),
			test.WithAuthentication(test.Authenticated, cookie),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// form assertions
		form := test.NewForm(doc, "project-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("", title.Value)
		assert.Equal("Title (required)", title.Label)

		description := form.MustGetFieldByID("description")
		assert.Equal("", description.Value)
		assert.Equal("Description", description.Label)

		assert.Equal("New project", test.FindByText(doc, "button", "New project").Text())
	})

	t.Run("new project returns form with errors", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Post),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// form assertions
		form := test.NewForm(doc, "project-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("Title (required)", title.Label)
		assert.Equal("", title.Value)
		assert.Equal("cannot be blank", title.Error)

		description := form.MustGetFieldByID("description")
		assert.Equal("Description", description.Label)
		assert.Equal("", description.Value)
		assert.Equal("", description.Error)

		assert.Equal("New project", test.FindByText(doc, "button", "New project").Text())
	})

	t.Run("new project redirects to '/projects'", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Post),
			test.WithFormValues(
				test.FormValue{
					Key:   "title",
					Value: "Project 5",
				},
				test.FormValue{
					Key:   "description",
					Value: "",
				},
			),
		)
		res := test.Do(req)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// redirection assertions
		assert.Equal("http://localhost:3000/projects", res.Header.Get("Hx-Redirect"))

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from projects where title = 'Project 5'")
		assert.Equal(1, count)
	})

	t.Run("edit project displays form and data", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects/1/edit")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Get),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// data assertions
		assert.Equal("Project 1", doc.Find("h1").Text())
		assert.Equal("Draft", doc.Find("div[id='status'] label").Text())

		// form assertions
		form := test.NewForm(doc, "project-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("Title (required)", title.Label)
		assert.Equal("Project 1", title.Value)

		description := form.MustGetFieldByID("description")
		assert.Equal("Description", description.Label)
		assert.Equal("", description.Value)

		assert.Equal(1, test.FindByText(doc, "button", "Save project").Size())
		assert.Equal(1, test.FindByText(doc, "button", "Delete project").Size())
	})

	t.Run("update project returns form", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects/1")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Patch),
			test.WithFormValues(
				test.FormValue{
					Key:   "title",
					Value: "Updated Project 1",
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

		// form assertions
		form := test.NewForm(doc, "project-form")

		title := form.MustGetFieldByID("title")
		assert.Equal("Title (required)", title.Label)
		assert.Equal("Updated Project 1", title.Value)

		description := form.MustGetFieldByID("description")
		assert.Equal("Description", description.Label)
		assert.Equal("", description.Value)

		assert.Equal(1, test.FindByText(doc, "button", "Save project").Size())
		assert.Equal(1, test.FindByText(doc, "button", "Delete project").Size())

		// oob assertions
		assert.Equal("Updated Project 1", doc.Find("h1").Text())

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from projects where title = 'Updated Project 1'")
		assert.Equal(1, count)
	})

	t.Run("toggle project returns switch", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects/1/toggle")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Patch),
		)
		res := test.Do(req)
		doc := test.Doc(res)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// switch assertions
		assert.Equal("Published", doc.Find("div[id='status'] label").Text())

		// db assertions
		var published bool
		handler.DB.Get(&published, "select projects.published from projects where id = 1")
		assert.True(published)
	})

	t.Run("delete project redirects to '/projects'", func(t *testing.T) {
		req := test.NewRequest(
			test.WithUrl(fmt.Sprintf("%s/%s", server.URL, "projects/1")),
			test.WithAuthentication(test.Authenticated, cookie),
			test.WithMethod(test.Delete),
		)
		res := test.Do(req)
		assert := assert.New(t)

		// status code assertions
		assert.Equal(200, res.StatusCode)

		// redirection assertions
		assert.Equal("http://localhost:3000/projects", res.Header.Get("Hx-Redirect"))

		// db assertions
		var count int
		handler.DB.Get(&count, "select count(*) from projects where id = 1")
		assert.Equal(0, count)
	})
}
