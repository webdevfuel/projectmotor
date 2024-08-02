package test

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

// A Form is a representation of a HTML form element, with a slice of fields.
type Form struct {
	Fields []*Field
}

// A Field is a representation of multiple HTML elements commonly used in the
// web application, including a label, input and span (error).
type Field struct {
	ID    string
	Name  string
	Value string
	Label string
	Error string
}

// NewForm returns a Form. It parses the provided goquery.Document and finds
// all fields within a form element with provided id.
func NewForm(doc *goquery.Document, fid string) *Form {
	fields := []*Field{}

	inputSelector := fmt.Sprintf(`form[id="%s"] input`, fid)
	textareaSelector := fmt.Sprintf(`form[id="%s"] textarea`, fid)

	doc.Find(inputSelector).Each(func(_ int, s *goquery.Selection) {
		id := s.AttrOr("id", "")
		name := s.AttrOr("name", "")
		value := s.AttrOr("value", "")

		lb := labelSelector(id)
		label := doc.Find(lb).Text()

		n := s.Next()

		fields = append(fields, &Field{
			ID:    id,
			Name:  name,
			Value: value,
			Label: label,
			Error: n.Text(),
		})
	})

	doc.Find(textareaSelector).Each(func(_ int, s *goquery.Selection) {
		id := s.AttrOr("id", "")
		name := s.AttrOr("name", "")
		value := s.Text()

		lb := labelSelector(id)
		label := doc.Find(lb).Text()

		n := s.Next()

		fields = append(fields, &Field{
			ID:    id,
			Name:  name,
			Value: value,
			Label: label,
			Error: n.Text(),
		})
	})

	return &Form{
		Fields: fields,
	}
}

func labelSelector(id string) string {
	return fmt.Sprintf(`label[for="%s"]`, id)
}

// GetFieldByID reports whether a field with the provided id exists in the
// Form and returns an empty Field if it doesn't exist.
func (form *Form) GetFieldByID(id string) (bool, *Field) {
	var exists bool
	var field *Field

	for _, f := range form.Fields {
		if f.ID == id {
			exists = true
			field = f
		}
	}

	return exists, field
}

// MustGetFieldByID returns a Field by the provided id and panics if it's not found.
func (form *Form) MustGetFieldByID(id string) *Field {
	for _, f := range form.Fields {
		if f.ID == id {
			return f
		}
	}

	log.Panic("field with id must be found")

	return &Field{}
}
