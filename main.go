package main

import (
	"log"
)

func main() {
	Initialize()
	defer Cleanup()

	err := LoadMaze("resources/maze.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}
	for {
		PrintMaze()

		input, err := ReadInput()
		if err != nil {
			log.Fatal("Failed to read input")
			break
		}

		MovePlayer(input)
		MoveGhosts()

		for _, ghost := range ghosts {
			if Player == *ghost {
				Lives = 0
			}
		}

		if input == "ESC" || Lives <= 0 || Dots == 0 {
			break
		}
	}
}
