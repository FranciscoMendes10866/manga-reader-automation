package types

type IManga struct {
	Name string
	URL  string
}

type INewMangaEntry struct {
	Thumbnail   string
	Name        string
	Description string
	Chapters    []IManga
}