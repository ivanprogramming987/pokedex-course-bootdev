package main

import (
	"bufio"
	"fmt"
	"github.com/ivanprogramming987/pokedex-course-bootdev/internal/pokecache"
	"os"
)

func main() {
	c := &config{
		supportedCommands: map[string]func(*config)error{
			"exit": commandExit,
			"help": commandHelp,
			"map":  commandMap,
			"mapb": commandMapb,
		},
		next: "https://pokeapi.co/api/v2/location-area/",
		previous: nil,
	}
	cache := pokecache.NewCache()
	scanner := bufio.NewScanner(os.Stdin)
	for ;; {
		// creates a command prompt
		fmt.Print("Pokedex > ")
		scanner.Scan()

		// gets the command you typed
		inputStr := cleanInput(scanner.Text())[0]

		command, ok := c.supportedCommands[inputStr]
		if !ok {
			fmt.Printf("Error: command %s not found\n", inputStr)
			continue
		}

		err := command(c)
		if err != nil {
			fmt.Errorf("Error: %w\n", err)
		}
	}
}
