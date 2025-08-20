package pokeapi

import (
    "net/http"
    "encoding/json"
    "io"
    "github.com/Ga11eaN/pokedex/internal/pokecache"
    "time"
    "fmt"
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

type PokemonAreaCall struct {
    Pokemons []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
    PokemonList struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"pokemon"`
}

func GetPokemons(url string) (PokemonAreaCall, error) {
    var bodyBytes []byte
    var err error

    bodyBytes, ok := cache.Get(url)
    if !ok {
        response, err := http.Get(url)
        if err != nil {
            // Server down / network issue
            return PokemonAreaCall{}, fmt.Errorf("could not connect: %v", err)
        }
        if response.StatusCode == 404 {
            // Bad location
            return PokemonAreaCall{}, fmt.Errorf("location area not found")
        }
        if response.StatusCode > 299 {
            // Other API/server error
            return PokemonAreaCall{}, fmt.Errorf("server error: %d", response.StatusCode)
        }
        bodyBytes, err = io.ReadAll(response.Body)
        if err != nil {
            return PokemonAreaCall{}, err
        }
        cache.Add(url, bodyBytes)
    }


    var pokemons PokemonAreaCall
    err = json.Unmarshal(bodyBytes, &pokemons)
    if err != nil {
        return PokemonAreaCall{}, err
    }

    return pokemons, nil
}

type Pokemon struct {
    Name string `json:"name"`
    BaseExperience int `json:"base_experience"`
}

func CatchPokemon(url string) (Pokemon, error) {
    response, err := http.Get(url)
    if err != nil {
        // Server down / network issue
        return Pokemon{}, fmt.Errorf("could not connect: %v", err)
    }
    if response.StatusCode == 404 {
        // Bad location
        return Pokemon{}, fmt.Errorf("Pokemon was not found")
    }
    if response.StatusCode > 299 {
        // Other API/server error
        return Pokemon{}, fmt.Errorf("server error: %d", response.StatusCode)
    }

    bodyBytes, err := io.ReadAll(response.Body)
    if err != nil {
        return Pokemon{}, err
    }

    var pokemon Pokemon
    err = json.Unmarshal(bodyBytes, &pokemon)
    if err != nil {
        return Pokemon{}, err
    }

    return pokemon, nil
}