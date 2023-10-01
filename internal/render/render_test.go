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
	// session.Put(r.Context(), "warning", "attention")
	// session.Put(r.Context(), "error", "you did it wrong")

	result := AddDefaultData(&td, r)
	if result == nil {
		t.Error("failed")
	}

	if result.Flash != "123" {
		t.Error("flash value 123 not found in session")
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
