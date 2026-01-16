package models

// Artist représente la structure d'un artiste telle qu'elle apparaît
// dans le tableau JSON de l'endpoint /api/artists
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	// Cette variable stocke l'URL vers l'objet "locations" (ex: /api/locations/1)
	Locations string `json:"locations"`
	// Ces variables stockent les autres URLs de l'API
	ConcertDates string `json:"concertDates"`
	Relations    string `json:"relations"`
}
