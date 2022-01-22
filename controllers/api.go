package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/hibiken/asynq"
)

func ScrapOnDemand(w http.ResponseWriter, r *http.Request) {
	var body types.IScrapOnDemandBody
	json.NewDecoder(r.Body).Decode(&body)

	newEntry := services.NewMangaEntry(body.URL)

	if len(newEntry.Name) > 2 {
		newMangaEntry := new(entities.MangaEntity)
		link, _ := services.UploadImageFromURL(newEntry.Thumbnail, newEntry.Name)
		newMangaEntry.Name = newEntry.Name
		newMangaEntry.Thumbnail = link
		newMangaEntry.Description = newEntry.Description

		config.Database.Create(&newMangaEntry)

		if len(newEntry.Chapters) > 0 && newMangaEntry.ID != "" {
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
			defer client.Close()

			task, _ := tasks.NewScrapSingleChapterTask(newMangaEntry.ID, newEntry.Chapters, newMangaEntry.Name)
			client.Enqueue(task)
		}
	}
}
