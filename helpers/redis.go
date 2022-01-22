package helpers

var RedisAddress = "localhost:6379"

var QueueProcessingConcurrency = 10

var QueuesDefinitions = map[string]int{
	"mangaScrap":     4,
	"mangaListScrap": 2,
	"chapterScrap":   4,
}
