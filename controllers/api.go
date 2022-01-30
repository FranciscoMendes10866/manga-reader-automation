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

	var databaseCategories []entities.CategoriesEntity
	config.Database.Table("manga_entities").Find(&databaseCategories)

	newEntry := services.NewMangaEntry(body.URL)

	if len(databaseCategories) > 0 {
		var categoriesToAdd []string
		for _, category := range newEntry.Categories {
			var categoryExists bool
			for _, databaseCategory := range databaseCategories {
				if category == databaseCategory.Name {
					categoryExists = true
				}
			}
			if !categoryExists {
				categoriesToAdd = append(categoriesToAdd, category)
			}
		}

		if len(categoriesToAdd) > 0 {
			for _, category := range categoriesToAdd {
				var newCategory entities.CategoriesEntity
				newCategory.Name = category
				config.Database.Create(&newCategory)
			}
		}
	}

	if len(newEntry.Name) > 2 {
		newMangaEntry := new(entities.MangaEntity)
		newMangaEntry.Name = newEntry.Name
		newMangaEntry.Thumbnail = newEntry.Thumbnail
		newMangaEntry.Description = newEntry.Description

		config.Database.Create(&newMangaEntry)

		if len(newEntry.Chapters) > 0 && newMangaEntry.ID != "" {
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
			defer client.Close()

			task, _ := tasks.NewScrapSingleChapterTask(newMangaEntry.ID, newEntry.Chapters, newMangaEntry.Name)
			client.Enqueue(task)

			var categoriesToAdd []string
			for _, category := range newEntry.Categories {
				var categoryExists bool
				for _, databaseCategory := range databaseCategories {
					if category == databaseCategory.Name {
						categoryExists = true
					}
				}
				if !categoryExists {
					categoriesToAdd = append(categoriesToAdd, category)
				}
			}

			if len(categoriesToAdd) > 0 {
				for _, category := range categoriesToAdd {
					var newCategory entities.CategoriesEntity
					config.Database.Where("name = ?", category).First(&newCategory)

					config.Database.Table("manga_categories").Create(map[string]interface{}{
						"manga_entity_id":      newMangaEntry.ID,
						"categories_entity_id": newCategory.ID,
					})
				}
			}
		}
	}
}
