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

	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Use(cors.Handler(helpers.CorsConfig))
	r.Use(controllers.AuthorizationMiddleware)
	r.Get("/api/manga/get-all", controllers.GetAllMangas)
	r.Get("/api/manga/get-details/{mangaId}", controllers.GetMangaDetails)
	r.Get("/api/chapter/{chapterId}", controllers.GetChapter)
	r.Get("/api/manga/latest", controllers.GetLatest)

	http.ListenAndServe(":3333", r)
}
