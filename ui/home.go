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
		return widget.NewLabel("Erreur de chargement")
	}

	list := container.NewVBox()

	for _, artist := range artists {
		a := artist

		btn := widget.NewButton(a.Name, func() {
			ShowArtistDetails(a)
		})

		list.Add(btn)
	}

	return container.NewVScroll(list)
}
