# Pokedex CLI

A command-line interface Pokedex built in Go that lets you explore the world of Pokemon using the PokeAPI.

## Features

- Map navigation to discover Pokemon locations
- Explore areas to find Pokemon
- Catch Pokemon with dynamic catch rates
- Inspect caught Pokemon's stats and details
- View your collected Pokemon in your Pokedex
- Caching system for faster response times

## Commands

- `help`: Displays available commands
- `exit`: Exit the Pokedex
- `map`: Display 20 location areas
- `mapb`: Display previous 20 location areas
- `explore <area-name>`: List Pokemon in a specific area
- `catch <pokemon-name>`: Try to catch a Pokemon
- `inspect <pokemon-name>`: View details of caught Pokemon
- `pokedex`: List all your caught Pokemon

## Technical Details

- Built in Go
- Uses PokeAPI for Pokemon data 
- Implements caching for improved performance
- Thread-safe operations

## Acknowledgments

This project was built following the excellent Go course curriculum from [Boot.dev](https://boot.dev). Their hands-on approach to learning Go through building real applications like this Pokedex made the learning experience both practical and enjoyable.