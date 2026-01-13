package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"groupietracker/geo"
	"groupietracker/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LocationsResponse struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

func ShowArtistMap(artist models.Artist) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Tournée de "+artist.Name,
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	cityList := container.NewVBox()
	cityScroll := container.NewVScroll(cityList)
	cityScroll.SetMinSize(fyne.NewSize(250, 450))

	loadLabel := widget.NewLabel("Initialisation...")
	loadProgress := widget.NewProgressBarInfinite()

	mapContainer := container.NewMax(container.NewVBox(loadLabel, loadProgress))

	split := container.NewHSplit(mapContainer, cityScroll)
	split.Offset = 0.7

	updateMap := func(newMap *canvas.Image) {
		if newMap.Resource != nil {
			mapContainer.Objects = []fyne.CanvasObject{newMap}
			mapContainer.Refresh()
			split.Refresh()
		}
	}

	go func() {
		resp, err := http.Get(artist.Locations)
		if err != nil {
			mapContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Erreur API")}
			return
		}
		defer resp.Body.Close()

		var locRes LocationsResponse
		json.NewDecoder(resp.Body).Decode(&locRes)

		var finalCoords []geo.Location

		for _, raw := range locRes.Locations {
			clean := strings.ReplaceAll(strings.ReplaceAll(raw, "-", ", "), "_", " ")
			cityName := strings.Title(clean)

			btn := widget.NewButton("📍 "+cityName, nil)
			btn.Disable()
			cityList.Add(btn)
		}
		cityList.Refresh()

		for i, raw := range locRes.Locations {
			clean := strings.ReplaceAll(strings.ReplaceAll(raw, "-", ", "), "_", " ")
			loadLabel.SetText(fmt.Sprintf("GPS %d/%d : %s", i+1, len(locRes.Locations), clean))

			loc, err := geo.GetCoordinates(clean)
			if err == nil {
				finalCoords = append(finalCoords, loc)

				if i < len(cityList.Objects) {
					currentLoc := loc
					cityList.Objects[i].(*widget.Button).OnTapped = func() {
						updateMap(geo.GenerateSingleCityMap(currentLoc))
					}
					cityList.Objects[i].(*widget.Button).Enable()
				}
			}
			time.Sleep(1100 * time.Millisecond)
		}

		viewAllBtn := widget.NewButtonWithIcon("VOIR TOUT", nil, func() {
			updateMap(geo.GenerateMapWithMultipleMarkers(finalCoords))
		})
		cityList.Objects = append([]fyne.CanvasObject{viewAllBtn, widget.NewSeparator()}, cityList.Objects...)
		cityList.Refresh()

		if len(finalCoords) > 0 {
			updateMap(geo.GenerateMapWithMultipleMarkers(finalCoords))
		} else {
			mapContainer.Objects = []fyne.CanvasObject{widget.NewLabel("Aucune ville trouvée.")}
		}
	}()

	return container.NewBorder(title, nil, nil, nil, split)
}
