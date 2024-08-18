package test

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

// A Form is a representation of a HTML form element, with a slice of fields.
type Form struct {
	Fields  []*Field
	Selects []*Select
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

// A Select is a representation of select HTML element, including a label and span (error).
type Select struct {
	ID      string
	Name    string
	Options []*Option
	Label   string
	Error   string
}

// A Option is a representation of option HTML element.
type Option struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
}

// NewForm returns a Form. It parses the provided goquery.Document and finds
// all fields within a form element with provided id.
func NewForm(doc *goquery.Document, fid string) *Form {
	fields := []*Field{}

	selects := []*Select{}

	doc.Find(fmt.Sprintf(`form[id="%s"] input`, fid)).Each(func(_ int, s *goquery.Selection) {
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

	doc.Find(fmt.Sprintf(`form[id="%s"] textarea`, fid)).Each(func(_ int, s *goquery.Selection) {
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

	doc.Find(fmt.Sprintf(`form[id="%s"] select`, fid)).Each(func(_ int, s *goquery.Selection) {
		id := s.AttrOr("id", "")
		name := s.AttrOr("name", "")

		lb := labelSelector(id)
		label := doc.Find(lb).Text()

		n := s.Next()

		options := []*Option{}

		s.Find("option").Each(func(_ int, s *goquery.Selection) {
			value := s.AttrOr("value", "")
			label := s.Text()
			_, selected := s.Attr("selected")
			_, disabled := s.Attr("disabled")

			options = append(options, &Option{
				Value:    value,
				Label:    label,
				Selected: selected,
				Disabled: disabled,
			})
		})

		selects = append(selects, &Select{
			ID:      id,
			Name:    name,
			Label:   label,
			Options: options,
			Error:   n.Text(),
		})
	})

	return &Form{
		Fields:  fields,
		Selects: selects,
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
	var f *Field

	for _, cf := range form.Fields {
		if cf.ID == id {
			f = cf
			break
		}
	}

	if f == nil {
		log.Panic("field with id must be found")
	}

	return f
}

// MustGetSelectByID returns a Select by the provided id and panics if it's not found.
func (form *Form) MustGetSelectByID(id string) *Select {
	var s *Select

	for _, cs := range form.Selects {
		if cs.ID == id {
			s = cs
			break
		}
	}

	if s == nil {
		log.Panic("field with id must be found")
	}

	return s
}
