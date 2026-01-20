//fichier qui permet de lancer l'application Fyne

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func StartApp() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	w.SetContent(BuildHome())
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
