package ui

import (
	"groupietracker/models"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowArtistDetails(artist models.Artist) {
	w := fyne.CurrentApp().NewWindow("Détails : " + artist.Name)

	members := "Membres:\n"
	for _, m := range artist.Members {
		members += "- " + m + "\n"
	}

	mapButton := widget.NewButton("Voir les concerts sur la carte", func() {
		mapWin := fyne.CurrentApp().NewWindow("Carte - " + artist.Name)
		mapWin.SetContent(ShowArtistMap(artist))
		mapWin.Resize(fyne.NewSize(850, 600))
		mapWin.Show()
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Année de création : "+strconv.Itoa(artist.CreationDate)),
		widget.NewLabel("Premier album : "+artist.FirstAlbum),
		widget.NewLabel(members),
		mapButton,
		widget.NewSeparator(),
		widget.NewButton("Fermer", func() {
			w.Close()
		}),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(450, 550))
	w.Show()
}
