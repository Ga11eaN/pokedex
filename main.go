package main

import (
    "fmt"
    "strings"
    "bufio"
    "os"
    "github.com/Ga11eaN/pokedex/internal/pokeapi"
)


type cliCommand struct {
    name string
    description string
    callback func(*Config) error
}

var commands map[string]cliCommand

type Config struct {
    Next string
    Previous string
}

var config Config

func main() {
    scanner := bufio.NewScanner(os.Stdin)
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
            "mapb" : {
                name: "mapb",
                description: "Show 20 previous areas on map",
                callback: commandMapb,
            },
        }

        if len(cleanedStr) < 1 {
            fmt.Println("Empty input. Please enter valid command or run 'help'")
        } else {
            cmd, exists := commands[cleanedStr[0]]
            if exists {
                err := cmd.callback(&config)
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

func commandExit(config *Config) error {
    fmt.Println("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp(config *Config) error {
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

func commandMap(config *Config) error {
    var url string
    if config.Next != "" {
        url = config.Next
    } else {
        url = "https://pokeapi.co/api/v2/location-area/"
    }
    return printAreaLocations(config, url)
}

func commandMapb(config *Config) error {
    if config.Previous == "" {
        fmt.Println("you're on the first page")
    } else {
        return printAreaLocations(config, config.Previous)
    }
    return nil
}

