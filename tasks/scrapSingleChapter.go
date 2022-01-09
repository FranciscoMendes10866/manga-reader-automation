package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/hibiken/asynq"
)

const TypeScrapSingleChapter = "scrap:chapter"

type ScrapSingleChapterPayload struct {
	MangaID  uint
	Chapters []services.IManga
}

func NewScrapSingleChapterTask(id uint, chapters []services.IManga) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapSingleChapterPayload{MangaID: id, Chapters: chapters})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScrapSingleChapter, payload, asynq.Queue("chapterScrap")), nil
}

func HandleScrapSingleChapterTask(ctx context.Context, t *asynq.Task) error {
	var payload ScrapSingleChapterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	mangaId := payload.MangaID
	chaptersData := payload.Chapters

	var databaseChapters []entities.ChapterEntity
	config.Database.Where("manga_id = ?", mangaId).Find(&databaseChapters)

	var newChapters []services.IManga

	for _, chapter := range chaptersData {
		var found bool
		for _, databaseChapter := range databaseChapters {
			if chapter.Name == databaseChapter.Name {
				found = true
				break
			}
		}
		if !found {
			newChapters = append(newChapters, chapter)
		}
	}

	if len(newChapters) > 0 {
		for _, chapter := range newChapters {
			pages := services.GetChapterPages(chapter.URL)

			if len(pages) > 0 {
				newChapterEntry := new(entities.ChapterEntity)
				newChapterEntry.Name = chapter.Name

				pagesString := ""
				for _, page := range pages {
					pagesString += page + ","
				}
				pagesString = pagesString[:len(pagesString)-1]
				newChapterEntry.Pages = pagesString
				newChapterEntry.MangaID = mangaId

				config.Database.Create(&newChapterEntry)
			}
		}
	}

	return nil
}
