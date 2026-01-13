package geo

import "sync"

// Cache pour éviter de re-géocoder les mêmes adresses
var (
	geocodeCache = make(map[string]Location)
	cacheMutex   sync.RWMutex
)
