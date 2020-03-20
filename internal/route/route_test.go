package route

import (
	"fmt"
	"net/http"
	"testing"
)

type urlPattern struct {
	pattern string
	isValid bool
}

type urlTestCase struct {
	urlPattern string
	patterns   []urlPattern
	pathParams []string
}

func hasPathParam(pathParams []string, pathParam string) bool {
	for _, paramActual := range pathParams {
		if pathParam == paramActual {
			return true
		}
	}

	return false
}

func (u *urlTestCase) doTest(t *testing.T) {
	regexData, err := UrlRegexConverter(u.urlPattern)
	if err != nil {
		t.Error(err)
	}

	testErrors := make([]string, 0)
	for _, p := range u.patterns {
		if regexData.pattern.MatchString(p.pattern) != p.isValid {
			testErrors = append(testErrors, fmt.Sprintf("test for pattern '%s' failed '%s'", u.urlPattern, p.pattern))
		}
	}

	for _, paramExpected := range u.pathParams {
		if !hasPathParam(regexData.pathParamKeys, paramExpected) {
			testErrors = append(testErrors, fmt.Sprintf("expected param '%s' not found", paramExpected))
		}
	}

	for _, errMsg := range testErrors {
		fmt.Println(errMsg)
	}

	if len(testErrors) > 0 {
		t.Errorf("test failed for url pattern '%s'", u.urlPattern)
	}
}

func TestUrlRegexConverter_SimplePath(t *testing.T) {
	testCase := urlTestCase{
		urlPattern: "/foo",
		pathParams: []string{},
		patterns: []urlPattern{
			{pattern: "/foo", isValid: true},
			{pattern: "/foobar", isValid: false},
			{pattern: "/foo/bar", isValid: false},
			{pattern: "/foo/", isValid: false},
			{pattern: "/fo", isValid: false},
			{pattern: "foo", isValid: false},
			{pattern: "", isValid: false},
			{pattern: "/", isValid: false},
		},
	}

	testCase.doTest(t)
}

func TestUrlRegexConverter_OneVariablePath(t *testing.T) {
	testCase := urlTestCase{
		urlPattern: "/foo/{id}",
		pathParams: []string{"id"},
		patterns: []urlPattern{
			{pattern: "/foo/bar", isValid: true},
			{pattern: "/foo/1", isValid: true},
			{pattern: "/foobar", isValid: false},
			{pattern: "/foo/bar/", isValid: false},
			{pattern: "/foo/1/", isValid: false},
			{pattern: "/foo/bar/abc", isValid: false},
			{pattern: "/foo/", isValid: false},
			{pattern: "/fo", isValid: false},
			{pattern: "foo", isValid: false},
			{pattern: "", isValid: false},
			{pattern: "/", isValid: false},
		},
	}

	testCase.doTest(t)
}

func TestRoutePathParam(t *testing.T) {
	testCase := urlTestCase{
		urlPattern: "/foo/{id}",
	}

	regexData, _:= UrlRegexConverter(testCase.urlPattern)

	route := route{regexData:regexData}

	request, err := http.NewRequest(http.MethodGet, "http://localhost/foo/2", nil)
	if err != nil {
		t.Error(err)
	}

	request = route.withPathParams(request)
	vars := Vars(request)
	if len(vars) != 1 {
		t.Errorf("wrong number of keys, expected '%d', got '%d'", 1, len(vars))
	}

	if value, ok := vars["id"] ; ok {
		if value != "2" {
			t.Errorf("expected value '%s' for key 'id', got '%s'", "2", value)
		}
	} else {
		t.Errorf("key '%s' not found", "id")
	}
}

