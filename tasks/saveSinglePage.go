package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/hibiken/asynq"
)

const TypeSaveSinglePageChapter = "save:image"

type SaveSinglePageChapterPayload struct {
	ChapterID string
	URL       string
}

func NewSaveSinglePageChapterTask(ChapterId string, url string) (*asynq.Task, error) {
	payload, err := json.Marshal(SaveSinglePageChapterPayload{ChapterID: ChapterId, URL: url})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSaveSinglePageChapter, payload, asynq.Queue("savePage")), nil
}

func HandleSaveSinglePageChapterTask(ctx context.Context, t *asynq.Task) error {
	var payload SaveSinglePageChapterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	chapterID := payload.ChapterID
	url := payload.URL

	newImageEntry := new(entities.PagesEntity)
	newImageEntry.ChapterID = chapterID
	newImageEntry.URL = url

	config.Database.Create(&newImageEntry)

	return nil
}
