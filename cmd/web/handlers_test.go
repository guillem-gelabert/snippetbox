package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	st := "./ui/static"
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
