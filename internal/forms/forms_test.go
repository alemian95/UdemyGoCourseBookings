package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/wathever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/wathever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/wathever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows does not have required when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	has := form.Has("wathever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "a")

	form = New(postedValues)
	has = form.Has("a")
	if !has {
		t.Error("form shows has not field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)
	form.MinLength("wathever", 3)
	if form.Valid() {
		t.Error("form shows minlength for non existent field")
	}

	isError := form.Errors.Get("wathever")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "aa")
	form = New(postedValues)
	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("form does not show minlength for field with 2 characters")
	}

	postedValues = url.Values{}
	postedValues.Add("a", "aaaa")
	form = New(postedValues)
	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("form does show minlength for field with 4 characters")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("form shows valid mail for non existent field")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "me@here.com")
	form = New(postedValues)

	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}

	postedValues = url.Values{}
	postedValues.Add("email", "x")
	form = New(postedValues)

	form.IsEmail("email")
	if form.Valid() {
		t.Error("got a valid for an invalid email")
	}
}
