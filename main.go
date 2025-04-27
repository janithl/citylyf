package main

import (
	"github.com/janithl/citylyf/internal"
	"github.com/janithl/citylyf/internal/ui"
)

var simRunner *internal.SimRunner

func startGame(gamePath *string) {
	if simRunner != nil { // if a game is already running, end it
		simRunner.EndGame()
	}
	simRunner = &internal.SimRunner{}
	simRunner.NewGame(gamePath)
	go simRunner.RunGameLoop() // start the game loop in a separate goroutine
}

func main() {
	ui.RunGame(startGame)
	if simRunner != nil { // if a game is running, end it
		simRunner.EndGame()
	}
}
