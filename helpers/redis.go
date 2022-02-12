package helpers

var RedisAddress = "localhost:6379"

var QueueProcessingConcurrency = 10

var QueuesDefinitions = map[string]int{
	"mangaScrap":         2,
	"mangaListScrap":     2,
	"chaptersScrap":      2,
	"savePage":           2,
	"singleChapterScrap": 2,
}
