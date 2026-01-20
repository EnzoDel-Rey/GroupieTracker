// Ce fichier contient les informations affichées lorsque l'on selectionne un artiste dans l'appli.

package ui

import (
	"fmt"
	"groupietracker/models"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowArtistDetails(artist models.Artist) {
	// Image du groupe
	res, err := fyne.LoadResourceFromURLString(artist.Image)
	var img *canvas.Image
	if err == nil {
		img = canvas.NewImageFromResource(res)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(200, 200))
	}

	// Informations du groupe
	name := widget.NewLabelWithStyle(artist.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})
	creation := widget.NewLabel(fmt.Sprintf("📅 Créé en : %d", artist.CreationDate))
	album := widget.NewLabel(fmt.Sprintf("💿 Premier album : %s", artist.FirstAlbum))

	members := widget.NewLabel(fmt.Sprintf("👥 Membres : %s", strings.Join(artist.Members, ", ")))
	members.Wrapping = fyne.TextWrapWord

	// Bouton pour aller voir la carte interactive des concerts
	mapBtn := widget.NewButtonWithIcon("VOIR LA TOURNÉE SUR LA CARTE", nil, func() {
		MainWindow.SetContent(ShowArtistMap(artist))
	})
	mapBtn.Importance = widget.HighImportance

	// Bouton Retour
	backBtn := widget.NewButton("Retour à l'accueil", func() {
		MainWindow.SetContent(BuildHome())
	})

	content := container.NewVBox(
		backBtn,
		name,
	)

	if img != nil {
		content.Add(img)
	}

	content.Add(creation)
	content.Add(album)
	content.Add(members)
	content.Add(widget.NewSeparator())
	content.Add(mapBtn)

	MainWindow.SetContent(container.NewPadded(container.NewVScroll(content)))
}
