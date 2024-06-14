package test

import (
	"io"
	"net/http"
)

type Method int

const (
	Get Method = iota
	Head
	Post
	Put
	Delete
	Connect
	Options
	Trace
	Patch
)

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

type Authentication int

const (
	Authenticated Authentication = iota
	Unauthenticated
)

type TestRequest struct {
	URL            string
	Body           []byte
	Method         Method
	Authentication Authentication
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
	req, _ := http.NewRequest(tr.Method.Build(), tr.URL, nil)
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

// WithUrl returns a function that sets the body
// on a TestRequest.
func WithBody(body []byte) func(*TestRequest) {
	return func(r *TestRequest) {
		r.Body = body
	}
}

// WithUrl returns a function that sets the method
// on a TestRequest.
func WithMethod(method Method) func(*TestRequest) {
	return func(r *TestRequest) {
		r.Method = method
	}
}

// WithUrl returns a function that sets the authentication
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
