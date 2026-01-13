/* code pour intégrer l'API nominatim qui permet de convertir une adresse en coordonnées */

package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Location struct {
	City      string
	Country   string
	Address   string
	Latitude  float64
	Longitude float64
}

type GeocodingResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func geocodeAddress(address string) (float64, float64, error) {
	encodedAddress := url.QueryEscape(address)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", encodedAddress)

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var results []GeocodingResult
	json.NewDecoder(resp.Body).Decode(&results)

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("adresse non trouvée: %s", address)
	}

	var lat, lon float64
	fmt.Sscanf(results[0].Lat, "%f", &lat)
	fmt.Sscanf(results[0].Lon, "%f", &lon)

	return lat, lon, nil
}

func GetCoordinates(address string) (Location, error) {
	cacheMutex.RLock()
	if cached, exists := geocodeCache[address]; exists {
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	lat, lon, err := geocodeAddress(address)
	if err != nil {
		return Location{}, err
	}

	loc := Location{
		Address:   address,
		Latitude:  lat,
		Longitude: lon,
	}

	cacheMutex.Lock()
	geocodeCache[address] = loc
	cacheMutex.Unlock()

	return loc, nil
}
