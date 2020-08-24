package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Player   string `json:"player"`
	Ghost    string `json:"ghost"`
	Wall     string `json:"wall"`
	Dot      string `json:"dot"`
	Pill     string `json:"pill"`
	Death    string `json:"death"`
	Space    string `json:"space"`
	UseEmoji bool   `json:"use_emoji"`
}

var Cfg Config

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
			MovePlayer(inp)
		default:
		}

		MoveGhosts()

		for _, ghost := range ghosts {
			if Player.row == ghost.row && Player.col == ghost.col {
				Lives -= 1
				if Lives != 0 {
					moveCursor(Player.row, Player.col)
					fmt.Print(Cfg.Death)
					time.Sleep(2000 * time.Millisecond)
					moveCursor(len(maze)+2, 0)
					Player.row, Player.col = Player.startRow, Player.startCol
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
