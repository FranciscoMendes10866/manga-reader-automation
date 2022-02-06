package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/FranciscoMendes10866/queues/config"
	"github.com/FranciscoMendes10866/queues/entities"
	"github.com/go-chi/chi/v5"
)

func GetAllMangas(w http.ResponseWriter, r *http.Request) {
	var mangas []entities.MangaEntity

	config.Database.Table("manga_entities").Find(&mangas)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mangas)
}

func GetMangaDetails(w http.ResponseWriter, r *http.Request) {
	mangaId := chi.URLParam(r, "mangaId")

	var results entities.MangaEntity

	config.Database.Table("manga_entities").Where("id = ?", mangaId).First(&results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
