package helpers

var RedisAddress = "localhost:6379"

var QueueProcessingConcurrency = 1

var QueuesDefinitions = map[string]int{
	"mangaScrap":     3,
	"mangaListScrap": 2,
	"chapterScrap":   3,
	"savePage":       2,
}
