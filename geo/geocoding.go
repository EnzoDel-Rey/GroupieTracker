package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Location est la structure utilisée par tout ton projet
type Location struct {
	City      string
	Country   string
	Address   string
	Latitude  float64
	Longitude float64
}

// googleGeocodeResponse permet de lire le JSON de Google
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

// --- DÉCLARATION DES VARIABLES GLOBALES DU PACKAGE ---
// Ces variables doivent être définies ici pour être accessibles par GetCoordinates
// GetCoordinates est la fonction appelée par ui/artistmap.go
func GetCoordinates(address string) (Location, error) {
	// 1. Vérification sécurisée du cache
	cacheMutex.RLock()
	if cached, exists := geocodeCache[address]; exists {
		cacheMutex.RUnlock()
		return cached, nil
	}
	cacheMutex.RUnlock()

	// 2. Appel à l'API Google (googleApiKey est dans mapview.go)
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

	// 3. Création de l'objet de retour
	loc := Location{
		Address:   address,
		Latitude:  res.Results[0].Geometry.Location.Lat,
		Longitude: res.Results[0].Geometry.Location.Lng,
	}

	// 4. Enregistrement sécurisé dans le cache
	cacheMutex.Lock()
	geocodeCache[address] = loc
	cacheMutex.Unlock()

	return loc, nil
}
