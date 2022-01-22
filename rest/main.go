package main

import (
	"net/http"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/controllers"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	config.Connect()
	config.ConnectBucket()

	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Use(cors.Handler(helpers.CorsConfig))
	r.Use(controllers.AuthorizationMiddleware)
	r.Post("/api/v1/manga/scrap-on-demand", controllers.ScrapOnDemand)

	http.ListenAndServe(":3333", r)
}
