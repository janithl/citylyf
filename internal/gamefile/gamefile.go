package gamefile

import (
	"encoding/json"
	"log"
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

// Save saves the current game state to a file at the specified path
func Save(path string) {
	var f *os.File
	var err error
	var saveGameJSON []byte

	if f, err = os.Create(path); err != nil {
		log.Println(err)
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
		log.Println(err)
		return
	}

	if _, err := f.Write(saveGameJSON); err != nil {
		log.Println(err)
		return
	}
}

// Load loads the game state from a file at the specified path
// and initializes the simulation with the loaded data
func Load(path string) {
	var fileData []byte
	var err error
	if fileData, err = os.ReadFile(path); err != nil {
		log.Println(err)
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

// GetSavesDir returns the directory where the game saves are stored.
func GetSavesDir() string {
	savesdir := ""
	if homedir, err := os.UserHomeDir(); err == nil {
		savesdir = homedir + "/.citylyf/saves"
		if err := os.MkdirAll(savesdir, os.ModePerm); err != nil {
			// creating the saves directory failed
			log.Println(err)
			savesdir = ""
		}
	} else {
		log.Println(err)
	}

	return savesdir
}

// GetDirFiles returns a list of files in the specified directory.
func GetDirFiles(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return []string{}
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames
}
