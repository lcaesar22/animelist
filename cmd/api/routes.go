package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(http.HandlerFunc(app.notFoundResponse))
	r.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	r.Get("/v1/anime", app.requirePermission("animes:read", app.listAnimesHandler))
	r.Get("/v1/healthcheck", app.healthcheckHandler)
	r.Post("/v1/anime", app.requirePermission("animes:read", app.createAnimeHandler))
	r.Get("/v1/anime/{id}", app.requirePermission("animes:read", app.showAnimeHandler))
	r.Patch("/v1/anime/{id}", app.requirePermission("animes:write", app.updateAnimeHandler))
	r.Delete("/v1/anime/{id}", app.requirePermission("animes:write", app.deleteAnimeHandler))

	//  users route
	r.Post("/v1/users", app.registerUserHandler)
	r.Put("/v1/users/activated", app.activateUserHandler)
	r.Post("/v1/tokens/authentication", app.createAuthenticationHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(r))))
}
