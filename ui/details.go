package ui

import (
	"groupietracker/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowArtistDetails(artist models.Artist) {
	w := fyne.CurrentApp().NewWindow(artist.Name)

	members := "Membres:\n"
	for _, m := range artist.Members {
		members += "- " + m + "\n"
	}

	content := container.NewVBox(
		widget.NewLabel(artist.Name),
		widget.NewLabel("Création: " + string(rune(artist.CreationDate))),
		widget.NewLabel("Premier album: " + artist.FirstAlbum),
		widget.NewLabel(members),
		widget.NewButton("Fermer", func() {
			w.Close()
		}),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 400))
	w.Show()
}
