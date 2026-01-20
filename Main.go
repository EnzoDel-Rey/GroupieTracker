// Main.go qui lance l'application et charge les villes en arrière plan pour éviter
// un long temps de chargement des données géographiques

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"groupietracker/geo"
	"groupietracker/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	ui.MainWindow = myApp.NewWindow("Groupie Tracker")

	go PreloadAllLocations()

	ui.MainWindow.SetContent(ui.BuildHome())

	ui.MainWindow.ShowAndRun()
}

func PreloadAllLocations() {
	fmt.Println("--- DÉMARRAGE DU PRÉ-CHARGEMENT DES COORDONNÉES API ---")

	allAddresses := getAllUniqueAddressesFromAPI()
	if len(allAddresses) == 0 {
		fmt.Println("Aucune ville à pré-charger ou erreur API.")
		return
	}

	fmt.Printf("%d villes uniques à traiter.\n", len(allAddresses))

	for i, address := range allAddresses {

		cleanAddr := strings.ReplaceAll(strings.ReplaceAll(address, "-", ", "), "_", " ")

		fmt.Printf("[%d/%d] Géocodage API : %s\n", i+1, len(allAddresses), cleanAddr)

		_, err := geo.GetCoordinates(cleanAddr)
		if err != nil {
			fmt.Printf("Erreur géocodage pour %s: %v\n", cleanAddr, err)
		}

		time.Sleep(1100 * time.Millisecond)
	}

	fmt.Println("--- TOUTES LES VILLES SONT PRÊTES EN MÉMOIRE ! ---")
}

func getAllUniqueAddressesFromAPI() []string {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		fmt.Println("Erreur lors de l'accès à l'API Locations:", err)
		return nil
	}
	defer resp.Body.Close()

	var data struct {
		Index []struct {
			Locations []string `json:"locations"`
		} `json:"index"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Erreur de décodage JSON:", err)
		return nil
	}

	uniqueMap := make(map[string]bool)
	var uniqueList []string

	for _, entry := range data.Index {
		for _, city := range entry.Locations {
			if !uniqueMap[city] {
				uniqueMap[city] = true
				uniqueList = append(uniqueList, city)
			}
		}
	}

	return uniqueList
}
