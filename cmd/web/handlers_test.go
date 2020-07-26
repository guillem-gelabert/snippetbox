package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	st := "./../../ui/static"
	var static *string = &st

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes(static))
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("expected %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("Expected body to equal %q", "OK")
	}
}

func TestShowSnippet(t *testing.T) {
	st := "./ui/static"
	var static *string = &st

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes(static))
	defer ts.Close()

	testCases := []struct {
		desc         string
		urlPath      string
		expectedCode int
		expectedBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("An old silent pond...")},
		{"Non-Existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			code, _, body := ts.get(t, tC.urlPath)
			if code != tC.expectedCode {
				t.Errorf("expected %d to be %d", code, tC.expectedCode)
			}
			if !bytes.Contains(body, tC.expectedBody) {
				t.Errorf("expected %s to be %s", body, tC.expectedBody)
			}
		})
	}
}
