package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/hibiken/asynq"
)

const TypeScrapSingleChapter = "scrap:chapter:single"

type ScrapSingleChapterPayload struct {
	ChapterURL  string
	ChapterName string
	MangaID     string
}

func NewScrapSingleChapterTask(chapterUrl string, chapterName string, mangaId string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapSingleChapterPayload{ChapterURL: chapterUrl, ChapterName: chapterName, MangaID: mangaId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScrapSingleChapter, payload, asynq.Queue("singleChapterScrap")), nil
}

func HandleScrapSingleChapterTask(ctx context.Context, t *asynq.Task) error {
	var payload ScrapSingleChapterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	pages := services.GetChapterPages(payload.ChapterURL)

	newChapterEntry := new(entities.ChapterEntity)
	newChapterEntry.Name = payload.ChapterName
	newChapterEntry.MangaID = payload.MangaID

	config.Database.Create(&newChapterEntry)

	for _, page := range pages {
		client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
		defer client.Close()

		task, _ := NewSaveSinglePageChapterTask(newChapterEntry.ID, page)
		client.Enqueue(task)
	}

	return nil
}
