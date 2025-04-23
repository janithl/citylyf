package main

import (
	"github.com/janithl/citylyf/internal"
	"github.com/janithl/citylyf/internal/ui"
)

func main() {
	simRunner := internal.SimRunner{}
	simRunner.NewGame(nil)
	go simRunner.RunGameLoop()
	ui.RunGame()
	simRunner.EndGame()
}
