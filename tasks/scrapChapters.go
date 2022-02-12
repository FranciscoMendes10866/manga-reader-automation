package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/hibiken/asynq"
)

const TypeScrapChapters = "scrap:chapters"

type ScrapChapterPayload struct {
	MangaID   string
	MangaName string
	Chapters  []types.IManga
}

func NewScrapChaptersTask(id string, chapters []types.IManga, MangaName string) (*asynq.Task, error) {
	payload, err := json.Marshal(ScrapChapterPayload{MangaID: id, MangaName: MangaName, Chapters: chapters})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeScrapChapters, payload, asynq.Queue("chaptersScrap")), nil
}

func HandleScrapChaptersTask(ctx context.Context, t *asynq.Task) error {
	var payload ScrapChapterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	mangaId := payload.MangaID
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
			client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
			defer client.Close()

			task, _ := NewScrapSingleChapterTask(chapter.URL, chapter.Name, mangaId)
			client.Enqueue(task)
		}
	}

	return nil
}
