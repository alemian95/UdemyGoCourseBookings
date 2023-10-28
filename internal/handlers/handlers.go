package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alemian95/UdemyGoCourseBookings/internal/config"
	"github.com/alemian95/UdemyGoCourseBookings/internal/driver"
	"github.com/alemian95/UdemyGoCourseBookings/internal/forms"
	"github.com/alemian95/UdemyGoCourseBookings/internal/helpers"
	"github.com/alemian95/UdemyGoCourseBookings/internal/models"
	"github.com/alemian95/UdemyGoCourseBookings/internal/render"
	"github.com/alemian95/UdemyGoCourseBookings/internal/repository"
	"github.com/alemian95/UdemyGoCourseBookings/internal/repository/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl.html", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl.html", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl.html", &models.TemplateData{})
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl.html", &models.TemplateData{})
}

func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type AvailabilityJsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) JsonSearchAvailability(w http.ResponseWriter, r *http.Request) {
	response := AvailabilityJsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl.html", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) RoomsMajorsSuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors-suite.page.tmpl.html", &models.TemplateData{})
}

func (m *Repository) RoomsGeneralsQuarters(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals-quarters.page.tmpl.html", &models.TemplateData{})
}
