// Ce fichier appelle 2 API :
// 1- Une API Google Geocode qui convertit les adresses en coordonnées
// 2- Une API Google Maps qui génère une map monde où l'on peut mettre des curseurs dessus.

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

// Structure qui permet de lire la réponse en Json de l'API google qui convertit les adresses en coordonnées
type googleGeocodeResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

func GetCoordinates(address string) (Location, error) {
	cacheMutex.RLock()
	if cached, exists := geocodeCache[address]; exists {
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	// Appel à l'API Google
	apiURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s",
		url.QueryEscape(address), googleApiKey)

	resp, err := http.Get(apiURL)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	var res googleGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return Location{}, err
	}

	if res.Status != "OK" || len(res.Results) == 0 {
		return Location{}, fmt.Errorf("Google error: %s pour %s", res.Status, address)
	}

	loc := Location{
		Address:   address,
		Latitude:  res.Results[0].Geometry.Location.Lat,
		Longitude: res.Results[0].Geometry.Location.Lng,
	}

	cacheMutex.Lock()
	geocodeCache[address] = loc
	cacheMutex.Unlock()

	return loc, nil
}
