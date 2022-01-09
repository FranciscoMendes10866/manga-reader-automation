package main

import (
	"log"

	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/hibiken/asynq"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeScheduled, tasks.HandleScheduledTask)

	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisAddr}, nil)

	task, _ := tasks.NewScheduledTask()

	entryID, err := scheduler.Register("@every 10s", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	scheduler.Run()
}
