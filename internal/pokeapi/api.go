package pokeapi

import (
    "net/http"
    "encoding/json"
    "io"
    "github.com/Ga11eaN/pokedex/internal/pokecache"
    "time"
)

var cache pokecache.Cache

func init() {
    cache = pokecache.NewCache(30 * time.Second)
}

type LocationAreaCall struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func MapGet(url string) (LocationAreaCall, error) {
    var bodyBytes []byte
    var err error

    bodyBytes, ok := cache.Get(url)
    if !ok {
        response, err := http.Get(url)
        if err != nil {
            return LocationAreaCall{}, err
        }
        defer response.Body.Close()

        bodyBytes, err = io.ReadAll(response.Body)
        if err != nil {
            return LocationAreaCall{}, err
        }
        cache.Add(url, bodyBytes)
    }


    var locations LocationAreaCall
    err = json.Unmarshal(bodyBytes, &locations)
    if err != nil {
        return LocationAreaCall{}, err
    }

    return locations, nil
}

