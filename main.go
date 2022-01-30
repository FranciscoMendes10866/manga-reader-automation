package main

import (
	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	config.Connect()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: helpers.RedisAddress},
		asynq.Config{
			Concurrency: helpers.QueueProcessingConcurrency,
			Queues:      helpers.QueuesDefinitions,
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeScheduled, tasks.HandleScheduledTask)
	mux.HandleFunc(tasks.TypeScrapSingleManga, tasks.HandleScrapSingleMangaTask)
	mux.HandleFunc(tasks.TypeScrapSingleChapter, tasks.HandleScrapSingleChapterTask)

	srv.Run(mux)
}
