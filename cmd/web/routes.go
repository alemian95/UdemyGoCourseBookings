package main

import (
	"net/http"

	"github.com/alemian95/go-bookings/pkg/config"
	"github.com/alemian95/go-bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)

	mux.Get("/rooms/generals-quarters", handlers.Repo.RoomsGeneralsQuarters)
	mux.Get("/rooms/majors-suite", handlers.Repo.RoomsMajorsSuite)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
