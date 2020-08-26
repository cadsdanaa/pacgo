package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Player              string        `json:"player"`
	Wall                string        `json:"wall"`
	Dot                 string        `json:"dot"`
	Pill                string        `json:"pill"`
	Death               string        `json:"death"`
	Space               string        `json:"space"`
	Ghost               string        `json:"ghost"`
	BlueGhost           string        `json:"blue_ghost"`
	PillDurationSeconds time.Duration `json:"pill_duration"`
	UseEmoji            bool          `json:"use_emoji"`
}

var Cfg Config
var previousDirection string

func main() {
	Initialize()
	defer Cleanup()

	err := LoadMaze("resources/maze.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}
	err = loadConfig("resources/config.json")
	if err != nil {
		log.Fatal("Failed to load config")
		return
	}

	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := ReadInput()
			if err != nil {
				ch <- "ESC"
				log.Fatal("Failed to read input")
			}
			ch <- input
		}
	}(input)

	for {
		PrintMaze()

		select {
		case inp := <-input:
			if inp == "ESC" {
				Lives = 0
			}
			validMove := MovePlayer(inp)
			if validMove {
				previousDirection = inp
			}
		default:
			MovePlayer(previousDirection)
		}

		MoveGhosts()

		for _, ghost := range ghosts {
			if Player.row == ghost.position.row && Player.col == ghost.position.col {
				if ghost.status == Normal {
					handleNormalGhostContact()
				} else {
					handleBlueGhostContact(ghost)
				}
			}
		}

		if Lives == 0 || Dots == 0 {
			if Lives == 0 {
				moveCursor(Player.row, Player.col)
				fmt.Print(Cfg.Death)
				moveCursor(len(maze)+2, 0)
			}
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func handleNormalGhostContact() {
	Lives -= 1
	if Lives != 0 {
		moveCursor(Player.row, Player.col)
		fmt.Print(Cfg.Death)
		time.Sleep(2000 * time.Millisecond)
		moveCursor(len(maze)+2, 0)
		Player.row, Player.col = Player.startRow, Player.startCol
	}
}
func handleBlueGhostContact(ghost *ghost) {
	ghost.position.row, ghost.position.col = ghost.position.startRow, ghost.position.startCol
	ghost.status = Normal
	score += 100
}

func loadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Cfg)
	if err != nil {
		return err
	}
	return nil
}
