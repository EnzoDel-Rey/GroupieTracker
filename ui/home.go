package ui

import (
	"groupietracker/api"
	"groupietracker/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildHome() fyne.CanvasObject {
	artists, err := api.FetchArtists()
	if err != nil {
		return widget.NewLabel("Erreur de chargement")
	}

	list := container.NewVBox()

	updateList := func(filtered []models.Artist) {
		list.Objects = nil
		for _, artist := range filtered {
			a := artist
			btn := widget.NewButton(a.Name, func() {
				ShowArtistDetails(a)
			})
			list.Add(btn)
		}
		list.Refresh()
	}

	updateList(artists)

	search := widget.NewEntry()
	search.SetPlaceHolder("Rechercher un groupe...")

	search.OnChanged = func(text string) {
		filtered := FilterArtists(artists, text)
		updateList(filtered)
	}

	return container.NewBorder(
		search, nil, nil, nil,
		container.NewVScroll(list),
	)
}

var MainWindow fyne.Window

func OnArtistSelected(artist models.Artist) {

	mapView := ShowArtistMap(artist)

	MainWindow.SetContent(mapView)
}
