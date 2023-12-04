package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("form is invalid when it should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form is valid even when required fields are missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form is shown as invalid even if it has required fields")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.Has("a")
	if form.Valid() {
		t.Error("form is shown as valid even if it does not have any fields")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	form.Has("a")
	if !form.Valid() {
		t.Error("form is shown as invalid even if it has a field")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("a", 1)
	if form.Valid() {
		t.Error("form is shown as valid length wise even if it is empty")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have gotten an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "aaa")

	form = New(postedData)
	form.MinLength("a", 2)
	if !form.Valid() {
		t.Error("field is shown as not meeting the minimum requirements even if it should")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have gotten an error, but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	email := form.Get("email")
	form.IsEmail(email)
	if form.Valid() {
		t.Error("email is shown as valid even if it's empty")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "r@f.com")

	form = New(postedValues)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("email is shown as invalid when it should be valid")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "r")

	form = New(postedValues)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("email is shown as valid when it should be invalid")
	}

}
