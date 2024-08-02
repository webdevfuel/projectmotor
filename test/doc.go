package test

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Docs reads the body from the provided response
// and returns a goquery.Document.
func Doc(res *http.Response) *goquery.Document {
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	return doc
}
