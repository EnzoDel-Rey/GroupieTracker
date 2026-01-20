// Ce fichier contient les informations de la page principale de l'appli

package ui

import (
	"encoding/json"
	"fmt"
	"groupietracker/api"
	"net/http"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var MainWindow fyne.Window

type AllLocationsResponse struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

func BuildHome() fyne.CanvasObject {
	artists, _ := api.FetchArtists()
	list := container.NewVBox()

	currentOpts := FilterOptions{
		CreationYearMin: 1950,
		MembersCount:    make(map[int]bool),
	}

	// Barre de recherche principale
	var allArtistNames []string
	for _, a := range artists {
		allArtistNames = append(allArtistNames, a.Name)
	}

	searchBar := widget.NewSelectEntry(allArtistNames)
	searchBar.SetPlaceHolder("Nom du groupe...")

	// --- 2. CONFIGURATION RECHERCHE VILLES ---
	var allCities []string
	locEntry := widget.NewSelectEntry(nil)
	locEntry.SetPlaceHolder("Ville de concert...")

	// Chargement des villes
	go func() {
		resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
		if err == nil {
			var data AllLocationsResponse
			json.NewDecoder(resp.Body).Decode(&data)
			uniqueCities := make(map[string]bool)
			for _, item := range data.Index {
				ArtistLocationsMap[item.ID] = item.Locations
				for _, city := range item.Locations {
					uniqueCities[city] = true
				}
			}
			for city := range uniqueCities {
				allCities = append(allCities, city)
			}
			locEntry.SetOptions(allCities)
			resp.Body.Close()
		}
	}()

	updateList := func() {
		filtered := ApplyAllFilters(artists, currentOpts)
		list.Objects = nil
		for _, a := range filtered {
			artist := a
			btn := widget.NewButton(artist.Name, func() { ShowArtistDetails(artist) })
			list.Add(btn)
		}
		list.Refresh()
	}

	// Suggestions pendant la recherche
	searchBar.OnChanged = func(s string) {
		currentOpts.SearchText = s
		if s == "" {
			searchBar.SetOptions(allArtistNames)
		} else {
			var filtered []string
			for _, name := range allArtistNames {
				if strings.Contains(strings.ToLower(name), strings.ToLower(s)) {
					filtered = append(filtered, name)
				}
			}
			searchBar.SetOptions(filtered)
		}
		searchBar.Show()
		updateList()
	}

	locEntry.OnChanged = func(s string) {
		currentOpts.LocationSearch = s
		if s == "" {
			locEntry.SetOptions(allCities)
		} else {
			var filtered []string
			lowS := strings.ToLower(s)
			for _, city := range allCities {
				if strings.Contains(strings.ToLower(city), lowS) {
					filtered = append(filtered, city)
				}
			}
			locEntry.SetOptions(filtered)
		}
		locEntry.Show()
		updateList()
	}

	// Visuel et emplacement des filtres de recherche
	sliderLabel := widget.NewLabel("Créé après : 1950")
	creationSlider := widget.NewSlider(1950, 2024)
	creationSlider.OnChanged = func(v float64) {
		currentOpts.CreationYearMin = int(v)
		sliderLabel.SetText(fmt.Sprintf("Créé après : %d", int(v)))
		updateList()
	}

	years := []string{"Toutes les années"}
	for i := 2024; i >= 1950; i-- {
		years = append(years, strconv.Itoa(i))
	}
	albumSelect := widget.NewSelect(years, func(s string) {
		if s == "Toutes les années" {
			currentOpts.FirstAlbumYear = 0
		} else {
			val, _ := strconv.Atoi(s)
			currentOpts.FirstAlbumYear = val
		}
		updateList()
	})

	membersRow := container.NewHBox()
	for i := 1; i <= 8; i++ {
		n := i
		membersRow.Add(widget.NewCheck(strconv.Itoa(n), func(b bool) {
			currentOpts.MembersCount[n] = b
			if !b {
				delete(currentOpts.MembersCount, n)
			}
			updateList()
		}))
	}

	filterContent := container.NewVBox(
		sliderLabel, creationSlider,
		widget.NewLabel("Année 1er Album :"), albumSelect,
		widget.NewLabel("Nombre de membres :"), membersRow,
		widget.NewLabel("Ville de concert (Suggestions) :"), locEntry,
	)

	header := container.NewVBox(
		widget.NewLabelWithStyle("GROUPIE TRACKER", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		searchBar,
		widget.NewAccordion(widget.NewAccordionItem("Filtres Avancés", filterContent)),
		widget.NewSeparator(),
	)

	updateList()
	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(list))
}
