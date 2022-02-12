package controllers

import (
	"encoding/json"
	"net/http"

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
