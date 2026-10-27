package main

import (
	"github.com/ivanprogramming987/pokedex-course-bootdev/internal/pokecache"
)

type config struct {
        supportedCommands map[string]func(*config, string)error
        next              string
	previous          *string
	cache             pokecache.Cache
	pokedex           map[string][]byte
}
