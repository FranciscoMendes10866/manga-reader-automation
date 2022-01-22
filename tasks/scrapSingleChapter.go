package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/hibiken/asynq"
)

const TypeScrapSingleChapter = "scrap:chapter"

type ScrapSingleChapterPayload struct {
	MangaID   string
	MangaName string
	Chapters  []types.IManga
}

func NewScrapSingleChapterTask(id string, chapters []types.IManga, MangaName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapSingleChapterPayload{MangaID: id, MangaName: MangaName, Chapters: chapters})
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
	mangaName := payload.MangaName
	chaptersData := payload.Chapters

	var databaseChapters []entities.ChapterEntity
	config.Database.Where("manga_id = ?", mangaId).Find(&databaseChapters)

	var newChapters []types.IManga

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

				var uploadedImages []string

				for _, page := range pages {
					imageURL, _ := services.UploadImageFromURL(page, mangaName)
					uploadedImages = append(uploadedImages, imageURL)
				}

				pagesString := ""
				for _, page := range uploadedImages {
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
