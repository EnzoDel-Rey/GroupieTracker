package ui

import (
	"groupietracker/models"
	"strings"
)

func FilterArtists(artists []models.Artist, query string) []models.Artist {
	var result []models.Artist

	query = strings.ToLower(query)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), query) {
			result = append(result, artist)
		}
	}

	return result
}
