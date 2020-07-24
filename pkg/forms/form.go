package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX is the compiled regexp for checking the format of an email address
var EmailRX = regexp.MustCompile("[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?")

// Form anonymously embeds a url.Values object (to hold the form data)
// and an Errors field to hold any validation errors  from the form data
type Form struct {
	url.Values
	Errors errors
}

// New initializes a custom Form struct from a PostForm/url.Values
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks that specific fields in the form data are present and not blank
// If any fields fail this check, adds tje appropriate message to the form errors
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks that a specific field in the form
// contains a maximum number of characters.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(
			field,
			fmt.Sprintf("This field is too long (maximum is %d characters)", d),
		)
	}
}

// PermittedValues checks that specific field in the form matches
// one of a set of specific permitted values.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// MinLength checks that a specific field in the form contains a minimum number of characters
func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MatchesPattern checks that a specific field matches a regular expression
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Valid returns true if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
