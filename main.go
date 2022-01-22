package main

import (
	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	config.Connect()
	config.ConnectBucket()

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				"mangaScrap":     4,
				"mangaListScrap": 2,
				"chapterScrap":   4,
			},
		},
	)

	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeScheduled, tasks.HandleScheduledTask)
	mux.HandleFunc(tasks.TypeScrapSingleManga, tasks.HandleScrapSingleMangaTask)
	mux.HandleFunc(tasks.TypeScrapSingleChapter, tasks.HandleScrapSingleChapterTask)

	srv.Run(mux)
}
