package api

import (
	"encoding/json"
	"groupietracker/models"
	"net/http"
)

// Collecte des données de l'API groupietracker

func FetchArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)

	return artists, err
}
