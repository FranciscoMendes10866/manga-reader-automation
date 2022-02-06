package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/hibiken/asynq"
)

const TypeScrapSingleManga = "scrap:manga"

type ScrapSingleMangaPayload struct {
	Manga types.IManga
}

func NewScrapSingleMangaTask(manga types.IManga) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapSingleMangaPayload{Manga: manga})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScrapSingleManga, payload, asynq.Queue("mangaScrap")), nil
}

func HandleScrapSingleMangaTask(ctx context.Context, t *asynq.Task) error {
	var payload ScrapSingleMangaPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	currentManga := payload.Manga

	var mangaInstance entities.MangaEntity
	config.Database.Where("name = ?", currentManga.Name).First(&mangaInstance)

	if mangaInstance.Name == "" && len(currentManga.Name) > 2 {
		var databaseCategories []entities.CategoriesEntity
		config.Database.Table("manga_categories").Find(&databaseCategories)

		newEntry := services.NewMangaEntry(currentManga.URL)

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
			newDatabaseEntry := new(entities.MangaEntity)
			newDatabaseEntry.Name = newEntry.Name
			newDatabaseEntry.Thumbnail = newEntry.Thumbnail
			newDatabaseEntry.Description = newEntry.Description

			config.Database.Create(&newDatabaseEntry)

			if len(newEntry.Chapters) > 0 && newDatabaseEntry.ID != "" {
				client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
				defer client.Close()

				task, _ := NewScrapSingleChapterTask(newDatabaseEntry.ID, newEntry.Chapters, newDatabaseEntry.Name)
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
							"manga_entity_id":      newDatabaseEntry.ID,
							"categories_entity_id": newCategory.ID,
						})
					}
				}
			}
		}
	} else if mangaInstance.Name != "" && len(currentManga.Name) > 2 {
		newChapters := services.GetMangaChapters(currentManga.URL)

		if len(newChapters) > 0 {
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
			defer client.Close()

			task, _ := NewScrapSingleChapterTask(mangaInstance.ID, newChapters, mangaInstance.Name)
			client.Enqueue(task)
		}
	}

	return nil
}
