package template

import "github.com/flosch/pongo2/v6"

var set = pongo2.NewSet("default", pongo2.MustNewLocalFileSystemLoader("./views"))

var (
	// Dashboard
	Dashboard *pongo2.Template = pongo2.Must(set.FromFile("dashboard/index.html"))
	// Login
	Login *pongo2.Template = pongo2.Must(set.FromFile("login/index.html"))
)
