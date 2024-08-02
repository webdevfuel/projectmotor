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

// FindByText returns a goquery.Selection with the provided selector,
// filtered by the provided text.
//
// It first finds all elements that match the selector, and then inside
// the Each function uses the Text method to compare against provided text.
func FindByText(doc *goquery.Document, selector string, text string) *goquery.Selection {
	var ts *goquery.Selection
	doc.Find(selector).Each(func(_ int, cs *goquery.Selection) {
		if cs.Text() == text {
			ts = cs
		}
	})
	return ts
}
