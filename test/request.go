package test

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form"
)

var encoder *form.Encoder

// Method is a int representation of HTTP methods, starting with Get at 0 and ending with Patch at 8.
type Method int

const (
	// GET Method
	Get Method = iota
	// HEAD Method
	Head
	// POST Method
	Post
	// PUT Method
	Put
	// DELETE Method
	Delete
	// CONNECT Method
	Connect
	// OPTIONS Method
	Options
	// TRACE Method
	Trace
	// PATCH Method
	Patch
)

// Build returns a string matching the HTTP method.
//
// An example:
//
//	log.Println(method.Get.Build())
//	"GET"
func (method Method) Build() string {
	switch method {
	case Get:
		return "GET"
	case Head:
		return "HEAD"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	case Connect:
		return "CONNECT"
	case Options:
		return "OPTIONS"
	case Trace:
		return "TRACE"
	case Patch:
		return "PATCH"
	}
	return ""
}

// Method is a int representation of authentication states for tests (authenticated and unauthenticated).
type Authentication int

const (
	// Authenticated (logged in) state for tests
	Authenticated Authentication = iota
	// Unauthenticated (logged out) state for tests
	Unauthenticated
)

// TestRequest is a representation of different options of a test request, used to build a http.Request from it.
type TestRequest struct {
	URL            string
	Body           []byte
	Method         Method
	Authentication Authentication
	UrlValues      url.Values
	IsForm         bool
	Header         http.Header
}

// NewRequest returns a http.Request build from provided TestRequest options.
//
// An example of usage inside a test:
//
//	request := test.NewRequest(
//	    test.WithUrl("/"),
//	    test.WithAuthentication(test.Authenticated, "")
//	)
func NewRequest(options ...func(*TestRequest)) *http.Request {
	tr := &TestRequest{
		URL:            "",
		Body:           []byte{},
		Method:         Get,
		Authentication: Unauthenticated,
		Header:         map[string][]string{},
	}
	for _, o := range options {
		o(tr)
	}
	req := &http.Request{}
	if tr.IsForm {
		req, _ = http.NewRequest(tr.Method.Build(), tr.URL, strings.NewReader(tr.UrlValues.Encode()))
	} else {
		req, _ = http.NewRequest(tr.Method.Build(), tr.URL, nil)
	}
	req.Header = tr.Header
	return req
}

// WithUrl returns a function that sets the url
// on a TestRequest.
func WithUrl(url string) func(*TestRequest) {
	return func(r *TestRequest) {
		r.URL = url
	}
}

// WithBody returns a function that sets the body
// on a TestRequest.
func WithBody(body []byte) func(*TestRequest) {
	return func(r *TestRequest) {
		r.Body = body
	}
}

// WithMethod returns a function that sets the method
// on a TestRequest.
func WithMethod(method Method) func(*TestRequest) {
	return func(r *TestRequest) {
		r.Method = method
	}
}

// WithAuthentication returns a function that sets the authentication
// and cookie/session values on a TestRequest.
func WithAuthentication(authentication Authentication, values ...string) func(*TestRequest) {
	return func(r *TestRequest) {
		r.Authentication = authentication

		switch authentication {
		case Authenticated:
			for _, v := range values {
				r.Header.Add("cookie", v)
			}
		case Unauthenticated:
			r.Header.Del("cookie")

		}
	}
}

// FormValue is a representation of a key-value pair
// for test requests.
type FormValue struct {
	Key   string
	Value string
}

// WithFormValues returns a function that sets the url values
// and content-type header as application/x-www-form-urlencoded
// on a TestRequest.
func WithFormValues(values ...FormValue) func(*TestRequest) {
	return func(r *TestRequest) {
		r.UrlValues = url.Values{}

		r.IsForm = true

		for _, v := range values {
			r.UrlValues.Set(v.Key, v.Value)
		}

		r.Header.Set("content-type", "application/x-www-form-urlencoded")
	}
}

// Do makes a request with the default http client and
// returns the response.
func Do(req *http.Request) *http.Response {
	response, _ := http.DefaultClient.Do(req)
	return response
}

// Body reads the body from the provided response
// and returns the string value.
func Body(res *http.Response) string {
	data, _ := io.ReadAll(res.Body)
	body := string(data)
	return body
}
