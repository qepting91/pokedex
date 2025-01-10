package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/qepting91/pokedex/internal/pokecache"
)

type config struct {
	nextUrl     string
	previousUrl string
	cache       *pokecache.Cache
	pokedex     map[string]pokemonResp
}

type locationAreaResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type pokemonEncounterResp struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemonResp struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View details about a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "View the pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandMap(cfg *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.nextUrl != "" {
		url = cfg.nextUrl
	}

	if data, ok := cfg.cache.Get(url); ok {
		var locationResp locationAreaResp
		err := json.Unmarshal(data, &locationResp)
		if err != nil {
			return err
		}

		cfg.nextUrl = ""
		if locationResp.Next != nil {
			cfg.nextUrl = *locationResp.Next
		}
		cfg.previousUrl = ""
		if locationResp.Previous != nil {
			cfg.previousUrl = *locationResp.Previous
		}

		for _, loc := range locationResp.Results {
			fmt.Println(loc.Name)
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg.cache.Add(url, body)

	var locationResp locationAreaResp
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}

	cfg.nextUrl = ""
	if locationResp.Next != nil {
		cfg.nextUrl = *locationResp.Next
	}
	cfg.previousUrl = ""
	if locationResp.Previous != nil {
		cfg.previousUrl = *locationResp.Previous
	}

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previousUrl == "" {
		fmt.Println("You're on the first page")
		return nil
	}

	url := cfg.previousUrl
	if data, ok := cfg.cache.Get(url); ok {
		var locationResp locationAreaResp
		err := json.Unmarshal(data, &locationResp)
		if err != nil {
			return err
		}

		cfg.nextUrl = ""
		if locationResp.Next != nil {
			cfg.nextUrl = *locationResp.Next
		}
		cfg.previousUrl = ""
		if locationResp.Previous != nil {
			cfg.previousUrl = *locationResp.Previous
		}

		for _, loc := range locationResp.Results {
			fmt.Println(loc.Name)
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg.cache.Add(url, body)

	var locationResp locationAreaResp
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}

	cfg.nextUrl = ""
	if locationResp.Next != nil {
		cfg.nextUrl = *locationResp.Next
	}
	cfg.previousUrl = ""
	if locationResp.Previous != nil {
		cfg.previousUrl = *locationResp.Previous
	}

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a location area name")
	}

	locationAreaName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationAreaName)

	if data, ok := cfg.cache.Get(url); ok {
		var encounterResp pokemonEncounterResp
		err := json.Unmarshal(data, &encounterResp)
		if err != nil {
			return err
		}
		fmt.Printf("Exploring %s...\n", locationAreaName)
		fmt.Println("Found Pokemon:")
		for _, encounter := range encounterResp.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg.cache.Add(url, body)

	var encounterResp pokemonEncounterResp
	err = json.Unmarshal(body, &encounterResp)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationAreaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounterResp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	pokemonName := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	if data, ok := cfg.cache.Get(url); ok {
		return handlePokemonCatch(cfg, data, pokemonName)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	cfg.cache.Add(url, body)
	return handlePokemonCatch(cfg, body, pokemonName)
}

func handlePokemonCatch(cfg *config, data []byte, pokemonName string) error {
	var pokemon pokemonResp
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}

	catchRate := 0.5 - float64(pokemon.BaseExperience)/1000.0
	if catchRate < 0.1 {
		catchRate = 0.1
	}

	randNum := rand.Float64()
	if randNum < catchRate {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.pokedex[pokemonName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandHelp(cfg *config, args ...string) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must provide a pokemon name")
	}

	pokemonName := args[0]
	pokemon, ok := cfg.pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}
func commandPokedex(cfg *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for name := range cfg.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
func main() {
	rand.Seed(time.Now().UnixNano())
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		cache:   pokecache.NewCache(5 * time.Minute),
		pokedex: make(map[string]pokemonResp),
	}
	commands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		command, exists := commands[commandName]

		if exists {
			err := command.callback(cfg, words[1:]...)
			if err != nil {
				fmt.Printf("Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.TrimSpace(lowered)
	words := strings.Fields(trimmed)
	return words
}
