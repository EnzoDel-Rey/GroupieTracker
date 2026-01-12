package ui

import (
	"groupietracker/api"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildHome() fyne.CanvasObject {
	artists, err := api.FetchArtists()

	if err != nil {
		return container.NewVBox(
			widget.NewLabel("Erreur de chargement de l'API"),
		)
	}

	list := container.NewVBox()

	for _, artist := range artists {
		list.Add(widget.NewLabel(artist.Name))
	}

	return container.NewVScroll(list)
}
