package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "github.com/Ga11eaN/pokedex/internal/pokeapi"
    "math/rand"
)


type cliCommand struct {
    name string
    description string
    callback func(cfg *Config, args ...string) error
}

var commands map[string]cliCommand

type Config struct {
    Next string
    Previous string
}

var config Config

var pokedex map[string]pokeapi.Pokemon

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    pokedex = make(map[string]pokeapi.Pokemon)
    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        text := scanner.Text()
        cleanedStr := cleanInput(text)


        commands = map[string]cliCommand{
            "exit": {
                name: "exit",
                description: "Exit the Pokedex",
                callback: commandExit,
            },
            "help": {
                name: "help",
                description: "Show available commands",
                callback: commandHelp,
            },
            "map": {
                name: "map",
                description: "Show next 20 areas on map",
                callback: commandMap,
            },
            "mapb": {
                name: "mapb",
                description: "Show 20 previous areas on map",
                callback: commandMapb,
            },
            "explore": {
                name: "explore {area}",
                description: "Explore area and search for pokemons",
                callback: commandExplore,
            },
            "catch": {
                name: "catch {pokemon}",
                description: "Catch pokemon",
                callback: commandCatch,
            },
        }

        if len(cleanedStr) < 1 {
            fmt.Println("Empty input. Please enter valid command or run 'help'")
        } else {
            cmd, exists := commands[cleanedStr[0]]
            if exists {
                err := cmd.callback(&config, cleanedStr[1:]...)
                if err != nil {
                    fmt.Printf("Error while running command, %v\n", err)
                }
            } else {
                fmt.Println("Unknown command")
            }
        }

    }

}


func cleanInput(text string) []string {
    cleaned := strings.ToLower(strings.TrimSpace(text))
    if cleaned == "" {
        return []string{}
    }
    return strings.Fields(cleaned)
}

func commandExit(config *Config, args ...string) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(config *Config, args ...string) error {
    fmt.Println("Welcome to the Pokedex!")
    fmt.Printf("Usage:\n\n")
    for _, cmd := range commands {
        fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }
    return nil
}

func printAreaLocations(config *Config, url string) error {
    resp, err := pokeapi.MapGet(url)
    if err != nil {
        fmt.Println(err)
    }

    for _, res := range(resp.Results) {
        fmt.Println(res.Name)
    }

    config.Next = resp.Next
    config.Previous = resp.Previous

    return nil
}

func commandMap(config *Config, args ...string) error {
    var url string
    if config.Next != "" {
        url = config.Next
    } else {
        url = "https://pokeapi.co/api/v2/location-area/"
    }
    return printAreaLocations(config, url)
}

func commandMapb(config *Config, args ...string) error {
    if config.Previous == "" {
        fmt.Println("you're on the first page")
    } else {
        return printAreaLocations(config, config.Previous)
    }
    return nil
}

func commandExplore(config *Config, args ...string) error {
    if len(args) < 0 {
        fmt.Println("Please enter location name to explore")
        return nil
    }
    location := args[0]
    fmt.Println("Exploring ", location)
    url := "https://pokeapi.co/api/v2/location-area/" + location +"/"
    return printPokemons(url)
}

func printPokemons(url string) error {

    resp, err := pokeapi.GetPokemons(url)
    if err != nil {
        fmt.Println(err)
    }
    for _, pokemon := range(resp.Pokemons) {
        fmt.Println(" -", pokemon.PokemonList.Name)
    }
    return nil
}

func commandCatch(config *Config, args ...string) error {
    if len(args) < 0 {
        fmt.Println("Please enter pokemon name to explore")
        return nil
    }
    pokemonName := args[0]
    fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
    url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

    resp, err := pokeapi.CatchPokemon(url)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    random_int := rand.Intn(700)
    if random_int > resp.BaseExperience {
        pokedex[pokemonName] = resp
        fmt.Println(pokemonName, "was caught!")
    } else {
        fmt.Println(pokemonName, "escaped!")
    }

    return nil
}
