package gamefile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/janithl/citylyf/internal/entities"
)

type SaveGame struct {
	Sim    *entities.Simulation
	LastID int
	Tiles  [][]entities.Tile
	Roads  []*entities.Road
}

func Save(path string) {
	var f *os.File
	var err error
	var saveGameJSON []byte

	if f, err = os.Create(path); err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	saveGame := SaveGame{
		Sim:    entities.Sim,
		LastID: entities.Sim.GetNextID(),
		Tiles:  entities.Sim.Geography.GetTiles(),
		Roads:  entities.Sim.Geography.GetRoads(),
	}

	if saveGameJSON, err = json.Marshal(saveGame); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := f.Write(saveGameJSON); err != nil {
		fmt.Println(err)
		return
	}
}

func Load(path string) {
	var fileData []byte
	var err error
	if fileData, err = os.ReadFile(path); err != nil {
		fmt.Println(err)
		return
	}

	jsonDecoder := json.NewDecoder(strings.NewReader(string(fileData)))
	saveGame := &SaveGame{}
	jsonDecoder.Decode(saveGame)

	entities.LoadSimulationFromSave(path, saveGame.Sim, uint32(saveGame.LastID), saveGame.Tiles, saveGame.Roads)
}

func CheckExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
