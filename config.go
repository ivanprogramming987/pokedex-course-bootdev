package main

type config struct {
        supportedCommands map[string]func(*config)error
        next              string
	previous          *string
}
