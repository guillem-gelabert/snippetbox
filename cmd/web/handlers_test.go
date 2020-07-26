package main

import (
	"bytes"
	"net/http"
	"net/url"
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

func TestSignupUser(t *testing.T) {
	st := "./../../ui/static"
	var static *string = &st

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes(static))
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")

	csrfToken := extractCSRFToken(t, body)

	testCases := []struct {
		desc         string
		userName     string
		userEmail    string
		userPAssword string
		csrfToken    string
		expectedCode int
		expectedBody []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty name", "Bob", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field is invalid")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field is too short (minimum is 10 characters)")},
		{"Duplicate email", "Bob", "dupe@example.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address is already in use")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tC.userName)
			form.Add("email", tC.userEmail)
			form.Add("password", tC.userPAssword)
			form.Add("csrf_token", tC.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			if code != tC.expectedCode {
				t.Errorf("expected %d to be %d", code, tC.expectedCode)
			}

			if !bytes.Contains(body, tC.expectedBody) {
				t.Errorf("expected %s to be %q", body, tC.expectedBody)
			}
		})
	}
}

func TestCreateSnippetForm(t *testing.T) {
	st := "./../../ui/static"
	var static *string = &st

	app := newTestApplication(t)
	ts := newTestServer(t, app.routes(static))
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create")

		actualHeader := headers.Get("Location")

		if http.StatusSeeOther != code {
			t.Errorf("expected %d to be %d", code, http.StatusSeeOther)
		}

		if actualHeader != "/user/login" {
			t.Errorf("expected %s to be %s", actualHeader, "/user/login")
		}
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "123123123")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		code, _, body := ts.get(t, "/snippet/create")

		if code != http.StatusOK {
			t.Errorf("expected %d to be %d", code, http.StatusOK)
		}

		if !bytes.Contains(body, []byte("<form action='/snippet/create' method='POST'>")) {
			t.Errorf("expected %s to be %q", body, "<form action='/snippet/create' method='POST'>")
		}
	})
}
