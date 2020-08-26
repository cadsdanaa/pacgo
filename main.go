package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	Initialize()
	defer Cleanup()
	InitSoundBoard()

	err := LoadMaze("resources/maze.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}
	err = LoadMazeConfig("resources/config.json")
	if err != nil {
		log.Fatal("Failed to load config")
		return
	}

	input := ReadUserInput()

	for {
		PrintMaze()

		select {
		case inp := <-input:
			if inp == "ESC" {
				Lives = 0
			}
			MovePlayer(inp)
		default:
			MovePlayer(previousDirection)
		}

		MoveGhosts()

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
