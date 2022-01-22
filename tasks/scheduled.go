package tasks

import (
	"context"

	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/services"
	"github.com/FranciscoMendes10866/queues/types"
	"github.com/hibiken/asynq"
)

const TypeScheduled = "scheduled:task"

type ScheduledPayload struct {
	Mangas []types.IManga
}

func NewScheduledTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeScheduled, nil, asynq.Queue("mangaListScrap")), nil
}

func HandleScheduledTask(ctx context.Context, t *asynq.Task) error {
	mangas := services.GetMangasList()

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: helpers.RedisAddress})
	defer client.Close()

	for _, manga := range mangas {
		task, _ := NewScrapSingleMangaTask(manga)
		client.Enqueue(task)
	}

	return nil
}
