package controllers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/go-chi/chi/v5"
)

func GetAllMangas(w http.ResponseWriter, r *http.Request) {
	var results []entities.MangaEntity
	config.Database.Table("manga_entities").Preload("Categories").Find(&results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func GetMangaDetails(w http.ResponseWriter, r *http.Request) {
	mangaId := chi.URLParam(r, "mangaId")

	var results entities.MangaEntity
	config.Database.Table("manga_entities").Where("id = ?", mangaId).Preload("Categories").Preload("Chapters").First(&results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func GetChapter(w http.ResponseWriter, r *http.Request) {
	chapterId := chi.URLParam(r, "chapterId")

	var result entities.ChapterEntity
	config.Database.Table("chapter_entities").Where("id = ?", chapterId).Preload("Pages").First(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GetLatest(w http.ResponseWriter, r *http.Request) {
	var chaptersList []entities.ChapterEntity
	config.Database.Table("chapter_entities").Order("created_at DESC").Limit(10000).Find(&chaptersList)

	var mangaIdsToFetch []string
	for _, currentChapter := range chaptersList {
		var exists bool
		for _, mangaId := range mangaIdsToFetch {
			if currentChapter.MangaID == mangaId {
				exists = true
				break
			}
		}
		if !exists {
			mangaIdsToFetch = append(mangaIdsToFetch, currentChapter.MangaID)
		}
	}

	var mangaList []entities.MangaEntity
	config.Database.Table("manga_entities").Where("id IN (?)", mangaIdsToFetch).Find(&mangaList)

	var mapped []entities.MangaEntity
	for _, currentManga := range mangaList {
		for _, currentChapter := range chaptersList {
			if currentManga.ID == currentChapter.MangaID {
				currentManga.Chapters = append(currentManga.Chapters, currentChapter)
			}
		}
		mapped = append(mapped, currentManga)
	}

	var remapped []entities.MangaEntity
	for _, currentManga := range mapped {
		sort.Slice(currentManga.Chapters, func(i, j int) bool {
			return currentManga.Chapters[i].CreatedAt.After(currentManga.Chapters[j].CreatedAt)
		})
		currentManga.Chapters = currentManga.Chapters[:2]
		remapped = append(remapped, currentManga)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(remapped)
}
