package main

import (
	"log"
	"time"
)

func main() {
	Initialize()
	defer Cleanup()

	err := LoadMaze("resources/maze.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}

	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := ReadInput()
			if err != nil {
				log.Fatal("Failed to read input")
				ch <- "ESC"
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
			if Player == *ghost {
				Lives = 0
			}
		}

		if Lives <= 0 || Dots == 0 {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}
