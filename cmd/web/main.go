package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alemian95/UdemyGoCourseBookings/internal/config"
	"github.com/alemian95/UdemyGoCourseBookings/internal/driver"
	"github.com/alemian95/UdemyGoCourseBookings/internal/handlers"
	"github.com/alemian95/UdemyGoCourseBookings/internal/helpers"
	"github.com/alemian95/UdemyGoCourseBookings/internal/models"
	"github.com/alemian95/UdemyGoCourseBookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const port = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {

	// what am i going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change thos to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 database=django_template_db user=example password=example")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	err = db.SQL.Ping()
	if err != nil {
		log.Fatal("Cannot ping database! Dying...")
	}
	log.Println("Database pinged")

	tc, cacheErr := render.CreateTemplateCache()
	if cacheErr != nil {
		return nil, cacheErr
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
