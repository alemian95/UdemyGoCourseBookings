package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alemian95/UdemyGoCourseBookings/internal/config"
	"github.com/alemian95/UdemyGoCourseBookings/internal/handlers"
	"github.com/alemian95/UdemyGoCourseBookings/internal/models"
	"github.com/alemian95/UdemyGoCourseBookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const port = ":8081"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// what am i going to put in the session
	gob.Register(models.Reservation{})

	// change thos to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, cacheErr := render.CreateTemplateCache()
	if cacheErr != nil {
		log.Fatal(cacheErr)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	log.Fatal(err)
}
