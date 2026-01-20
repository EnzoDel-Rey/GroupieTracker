// Ce fichier contient les filtres de recherche de l'appli.

package ui

import (
	"fmt"
	"groupietracker/models"
	"strings"
)

var ArtistLocationsMap = make(map[int][]string)

type FilterOptions struct {
	SearchText      string
	CreationYearMin int
	FirstAlbumYear  int
	MembersCount    map[int]bool
	LocationSearch  string
}

func ApplyAllFilters(artists []models.Artist, opt FilterOptions) []models.Artist {
	var filtered []models.Artist

	for _, a := range artists {
		// Filtre Nom
		if opt.SearchText != "" && !strings.Contains(strings.ToLower(a.Name), strings.ToLower(opt.SearchText)) {
			continue
		}
		// Filtre Année Création (Curseur - / +)
		if a.CreationDate < opt.CreationYearMin {
			continue
		}
		// Filtre Premier Album (Bouton déroulant avec une liste de dates)
		if opt.FirstAlbumYear != 0 {
			parts := strings.Split(a.FirstAlbum, "-")
			if len(parts) == 3 {
				var year int
				fmt.Sscanf(parts[2], "%d", &year)
				if year != opt.FirstAlbumYear {
					continue
				}
			}
		}
		// Filtre Membres (check boxes)
		if len(opt.MembersCount) > 0 {
			if !opt.MembersCount[len(a.Members)] {
				continue
			}
		}
		// Filtre Ville (string)
		if opt.LocationSearch != "" {
			cities, exists := ArtistLocationsMap[a.ID]
			if !exists {
				continue
			}
			found := false
			searchTerm := strings.ToLower(opt.LocationSearch)
			for _, c := range cities {
				if strings.Contains(strings.ToLower(c), searchTerm) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filtered = append(filtered, a)
	}
	return filtered
}
