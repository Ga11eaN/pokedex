package pokeapi_test


import (
    "testing"
    "github.com/Ga11eaN/pokedex/internal/pokeapi"
)

func TestMapGet(t *testing.T) {
    testURL := "https://pokeapi.co/api/v2/location-area/"
    resp, err := pokeapi.MapGet(testURL)
    if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Count == 0 {
		t.Error("expected count to be greater than 0")
	}
}

