package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func commandExit(c *config) error {
        fmt.Println("Closing the Pokedex... Goodbye!")
        os.Exit(0)
        return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	fmt.Println("exit: Use to exit the Pokedex")
	fmt.Println("help: Use to display all commands and what they do")
	fmt.Println("map: Displays the next 20 locations")
	fmt.Println("mapb: Displays the previous 20 locations")
	return nil
}

func commandMap(c *config) error {
	res, err := http.Get(c.next)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var area locationArea

	err = json.Unmarshal(body, &area)
	if err != nil {
		return err
	}

	for _, result := range area.Results {
		fmt.Println(result.Name)
	}

	c.next = area.Next
	c.previous = area.Previous

	return nil
}

func commandMapb(c *config) error {
	if c.previous == nil {
		fmt.Println("You are on first page.")
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
        var area locationArea

        err = json.Unmarshal(body, &area)
        if err != nil {
              	return err
        }

        for _, result := range area.Results {
                fmt.Println(result.Name)
	}

        c.next = area.Next
        c.previous = area.Previous

        return nil
}

type locationArea struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
