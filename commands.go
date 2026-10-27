package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

func commandExit(c *config, p string) error {
        fmt.Println("Closing the Pokedex... Goodbye!")
        os.Exit(0)
        return nil
}

func commandHelp(c *config, p string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("exit: Use to exit the Pokedex.")
	fmt.Println("help: Use to display all commands and what they do.")
	fmt.Println("map: Displays the next 20 locations.")
	fmt.Println("mapb: Displays the previous 20 locations.")
	fmt.Println("explore <area-name>: Displays all possible Pokemon encounters within the area <area-name>.")
	fmt.Println("catch <pokemon>: Attempt to catch <pokemon>. If success, adds <pokemon> to your Pokedex.")
	fmt.Println("inspect <pokemon>: Prints some useful attributes about <pokemon> if <pokemon> is in your Pokedex.")
	fmt.Println("pokedex: Shows all the Pokemon you have caught.")
	return nil
}

func commandMap(c *config, p string) error {
	cached_val, ok := c.cache.Get(c.next)
	if ok {
		var areas mapResponse

		err := json.Unmarshal(cached_val, &areas)
        	if err != nil {
                	return err
        	}

        	for _, result := range areas.Results {
                	fmt.Println(result.Name)
        	}

        	c.next = areas.Next
        	c.previous = areas.Previous

        	return nil
	}

	res, err := http.Get(c.next)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.cache.Add(c.next, body)

	var areas mapResponse

	err = json.Unmarshal(body, &areas)
	if err != nil {
		return err
	}

	for _, result := range areas.Results {
		fmt.Println(result.Name)
	}

	c.next = areas.Next
	c.previous = areas.Previous

	return nil
}

func commandMapb(c *config, p string) error {
	if c.previous == nil {
		fmt.Println("You are on first page.")
		return nil
	}

	cached_val, ok := c.cache.Get(*c.previous)
	if ok {
		var areas mapResponse

		err := json.Unmarshal(cached_val, &areas)
                if err != nil {
                        return err
                }

                for _, result := range areas.Results {
                        fmt.Println(result.Name)
                }

                c.next = areas.Next
                c.previous = areas.Previous

                return nil
	}

	res, err := http.Get(*c.previous)
        if err != nil {
              	return err
        }

        defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
       	if err != nil {
               	return err
       	}

	c.cache.Add(*c.previous, body)

        var areas mapResponse

        err = json.Unmarshal(body, &areas)
        if err != nil {
              	return err
        }

        for _, result := range areas.Results {
                fmt.Println(result.Name)
	}

        c.next = areas.Next
        c.previous = areas.Previous

        return nil
}

func commandExplore(c *config, p string) error {
	cached_val, ok := c.cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", p))

	if ok {
		var area exploreResponse

        	err := json.Unmarshal(cached_val, &area)
        	if err != nil {
                	return err
        	}

        	for _, result := range area.PokemonEncounters {
                	fmt.Printf("- %s\n", result.Pokemon.Name)
        	}

        	return nil
	}

	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", p))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.cache.Add(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", p), body)

	var area exploreResponse

	err = json.Unmarshal(body, &area)
        if err != nil {
                return err
        }

        for _, result := range area.PokemonEncounters {
                fmt.Printf("- %s\n", result.Pokemon.Name)
        }

	return nil
}

func commandCatch(c *config, p string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", p)

	cached_val, ok := c.cache.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", p))
	if ok {
		var pokemon catchResponse

        	err := json.Unmarshal(cached_val, &pokemon)
        	if err != nil {
                	return err
        	}

        	b := pokemon.BaseExperience

		acc := 0

		if b == 0 {
			fmt.Printf("Error: %s not found (if %s has multiple forms, specify the form (example: darmanitan-standard))\n", p, p)
			return nil
		} else if b > 0 && b < 50 {
			acc = 80
		} else if b >= 50 && b < 80 {
			acc = 60
		} else if b >= 80 && b < 120 {
			acc = 40
		} else {
			acc = 20
		}

        	r := rand.New(rand.NewSource(0))

        	if r.Intn(100) <= acc {
                	fmt.Printf("%s was caught!\n", p)
			fmt.Printf("You may now inspect it with the inspect command.\n")
                	c.pokedex[p] = cached_val
        	} else {
                	fmt.Printf("%s got away.\n", p)
        	}

        	return nil
	}

	res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", p))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	c.cache.Add(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", p), body)

	var pokemon catchResponse

	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	b := pokemon.BaseExperience

	acc := 0

	if b == 0 {
		fmt.Printf("Error: %s not found (if %s has multiple forms, specify the form (example: darmanitan-standard))\n", p, p)
		return nil
	} else if b > 0 && b < 50 {
                acc = 80
        } else if b >= 50 && b < 80 {
                acc = 60
        } else if b >= 80 && b < 120 {
                acc = 40
        } else {
                acc = 20
        }

	r := rand.New(rand.NewSource(1))

	if r.Intn(100) >= acc {
		fmt.Printf("%s was caught!\n", p)
		fmt.Printf("You may now inspect it with the inspect command.\n")
		c.pokedex[p] = body
	} else {
		fmt.Printf("%s got away.\n", p)
	}

	return nil
}

func commandInspect(c *config, p string) error {
	_, ok := c.pokedex[p]
	if !ok {
		fmt.Printf("You have not caught that Pokemon. (Run catch %s to attempt to catch it)\n", p)
		return nil
	}

	var pokemon catchResponse

	err := json.Unmarshal(c.pokedex[p], &pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Print("Stats:\n")
	fmt.Printf("  hp: %d\n", pokemon.Stats[0].BaseStat)
	fmt.Printf("  attack: %d\n", pokemon.Stats[1].BaseStat)
	fmt.Printf("  defense: %d\n", pokemon.Stats[2].BaseStat)
	fmt.Printf("  special attack: %d\n", pokemon.Stats[3].BaseStat)
	fmt.Printf("  special defense: %d\n", pokemon.Stats[4].BaseStat)
	fmt.Printf("  speed: %d\n", pokemon.Stats[5].BaseStat)
	fmt.Print("Types:\n")
	fmt.Printf("  %s\n", pokemon.Types[0].Type.Name)
	if len(pokemon.Types) > 1 {
		fmt.Printf("  %s\n", pokemon.Types[1].Type.Name)
	}

	return nil
}

func commandPokedex(c *config, p string) error {
	fmt.Printf("Your Pokedex:\n")
	for pokemon, _ := range c.pokedex {
		fmt.Printf(" - %s\n", pokemon)
	}

	return nil
}

type mapResponse struct {
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

type exploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type catchResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
