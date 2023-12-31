package render

import (
	"net/http"
	"testing"

	"github.com/alemian95/UdemyGoCourseBookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result == nil {
		t.Error("failed")
	}

	if result.Flash != "123" {
		t.Error("flash value 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww testWriter

	err = RenderTemplate(&ww, r, "home.page.tmpl.html", &models.TemplateData{})

	if err != nil {
		t.Error(err)
	}

	err = RenderTemplate(&ww, r, "wrong-template-name.page.tmpl.html", &models.TemplateData{})

	if err == nil {
		t.Error("rendered template that does not exists")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
