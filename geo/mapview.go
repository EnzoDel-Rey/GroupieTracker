package geo

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// clé API de google maps (sur compte de Flo)
const googleApiKey = "AIzaSyBnSLuZEXuqaEHcJC4K0Bozqs_yjkvxowI"

func GenerateMapWithMultipleMarkers(locations []Location) *canvas.Image {
	if len(locations) == 0 {
		return canvas.NewImageFromResource(nil)
	}

	baseURL := "https://maps.googleapis.com/maps/api/staticmap?"
	params := fmt.Sprintf("size=800x600&maptype=roadmap&language=fr&key=%s", googleApiKey)

	markers := ""
	for _, loc := range locations {
		markers += fmt.Sprintf("&markers=color:red|%.6f,%.6f", loc.Latitude, loc.Longitude)
	}

	return downloadGoogleMap(baseURL + params + markers)
}

func GenerateSingleCityMap(loc Location) *canvas.Image {
	mapURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/staticmap?center=%.6f,%.6f&zoom=4&size=800x600&maptype=roadmap&markers=color:blue|label:S|%.6f,%.6f&language=fr&key=%s",
		loc.Latitude, loc.Longitude, loc.Latitude, loc.Longitude, googleApiKey,
	)
	return downloadGoogleMap(mapURL)
}

func downloadGoogleMap(mapURL string) *canvas.Image {
	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest("GET", mapURL, nil)
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	resp, err := client.Do(req)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		return canvas.NewImageFromResource(nil)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	resource := fyne.NewStaticResource("google_map.png", data)
	img := canvas.NewImageFromResource(resource)

	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(600, 450))

	return img
}
