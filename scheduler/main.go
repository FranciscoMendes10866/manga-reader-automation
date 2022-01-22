package main

import (
	"log"

	"github.com/FranciscoMendes10866/queues/helpers"
	"github.com/FranciscoMendes10866/queues/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeScheduled, tasks.HandleScheduledTask)

	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: helpers.RedisAddress}, nil)

	task, _ := tasks.NewScheduledTask()

	entryID, err := scheduler.Register("@every 30m", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	scheduler.Run()
}
