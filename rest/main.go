package main

import (
	"encoding/json"
	"net/http"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	config.Connect()
	config.ConnectBucket()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Post("/api/v1/manga/scrap-on-demand", func(w http.ResponseWriter, r *http.Request) {
		req := map[string]string{}
		json.NewDecoder(r.Body).Decode(&req)

		newEntry := services.NewMangaEntry(req["url"])

		if len(newEntry.Name) > 2 {
			newMangaEntry := new(entities.MangaEntity)
			link, _ := services.UploadImageFromURL(newEntry.Thumbnail, newEntry.Name)
			newMangaEntry.Name = newEntry.Name
			newMangaEntry.Thumbnail = link
			newMangaEntry.Description = newEntry.Description

			config.Database.Create(&newMangaEntry)

			if len(newEntry.Chapters) > 0 && newMangaEntry.ID != 0 {
				client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
				defer client.Close()

				task, _ := tasks.NewScrapSingleChapterTask(newMangaEntry.ID, newEntry.Chapters, newMangaEntry.Name)
				client.Enqueue(task)
			}
		}
	})

	http.ListenAndServe(":3333", r)
}
