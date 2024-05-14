package validator

import (
	"fmt"
	"net/http"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-playground/form"
)

var decoder *form.Decoder

// A Validator is a struct of any type with a Validate method.
type Validator interface {
	Validate() error
}

// A Validated is a struct that holds the key, value and error strings.
// Key is used to grab a Validated instance inside templates, and should always be PascalCase.
// Value is used to pass to the client the previously available input value.
// Error is used to pass to the client a human-readable message.
type Validated struct {
	Key   string
	Value string
	Error string
}

// A ValidatedSlice is an alias to []Validated.
type ValidatedSlice []Validated

// GetByKey returns a Validated struct from a slice,
// and should always be PascalCase.
func (vs ValidatedSlice) GetByKey(key string) Validated {
	for _, v := range vs {
		if v.Key == key {
			return v
		}
	}
	return Validated{}
}

// NewValidatedSlice returns a ValidatedSlice.
func NewValidatedSlice() ValidatedSlice {
	return []Validated{}
}

// Validate reports whether the type implementing the Validator interface is
// valid, according to validation rules specified within Validate method.
//
// It also returns a ValidatedSlice with all information about the invalid
// type and the first error encountered while trying to validating.
//
// The error only occurs if there was an internal problem with the function,
// and bool should be used to track whether the validation was successful
// according to the rules.
func Validate(v Validator, r *http.Request) (bool, ValidatedSlice, error) {
	decoder = form.NewDecoder()
	err := r.ParseForm()
	if err != nil {
		return false, []Validated{}, err
	}
	err = decoder.Decode(&v, r.Form)
	if err != nil {
		return false, []Validated{}, err
	}
	err = v.Validate()
	if errors, ok := err.(validation.Errors); ok {
		parsedErrors := parseErrors(errors, v)
		return false, parsedErrors, nil
	}
	if err != nil {
		return false, []Validated{}, err
	}
	return true, []Validated{}, nil
}

func parseErrors(errors validation.Errors, data any) []Validated {
	emap := map[string]string{}
	vmap := []Validated{}
	for k, v := range errors {
		emap[k] = v.Error()
	}
	value := reflect.ValueOf(data).Elem()
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		kv := getKeyValue(value, fieldName, i)
		for _, kv := range kv {
			vmap = append(vmap, Validated{
				Key:   kv.Key,
				Value: kv.Value,
				Error: emap[kv.Key],
			})
		}
	}
	return vmap
}

type keyValue struct {
	Key   string
	Value string
}

func getKeyValue(val reflect.Value, fieldName string, i int) []keyValue {
	switch val.Field(i).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []keyValue{{
			Key:   fieldName,
			Value: fmt.Sprintf("%d", val.Field(i).Int()),
		}}
	case reflect.String:
		return []keyValue{{
			Key:   fieldName,
			Value: val.Field(i).String(),
		}}
	case reflect.Bool:
		return []keyValue{{
			Key:   fieldName,
			Value: fmt.Sprintf("%t", val.Field(i).Bool()),
		}}
	}
	return []keyValue{}
}
