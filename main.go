package main

import (
	"bufio"
	"fmt"
	"github.com/ivanprogramming987/pokedex-course-bootdev/internal/pokecache"
	"os"
	"time"
)

func main() {
	c := &config{
		supportedCommands: map[string]func(*config, string)error{
			"exit":    commandExit,
			"help":    commandHelp,
			"map":     commandMap,
			"mapb":    commandMapb,
			"explore": commandExplore,
			"catch":   commandCatch,
			"inspect": commandInspect,
			"pokedex": commandPokedex,
		},
		next: "https://pokeapi.co/api/v2/location-area/",
		previous: nil,
		cache: pokecache.NewCache(15 * time.Second),
		pokedex: map[string][]byte{},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		// creates a command prompt
		fmt.Print("Pokedex > ")
		scanner.Scan()

		// gets the command you typed
		cleanText := cleanInput(scanner.Text())
		inputCommand := cleanText[0]
		inputParameter := ""
		if len(cleanText) > 1 {
			inputParameter = cleanText[1]
		}

		command, ok := c.supportedCommands[inputCommand]
		if !ok {
			fmt.Printf("Error: command %s not found\n", inputCommand)
			continue
		}

		err := command(c, inputParameter)
		if err != nil {
			fmt.Errorf("Error: %w\n", err)
		}
	}
}
