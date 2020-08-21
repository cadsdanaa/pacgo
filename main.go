package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var maze []string

func main() {
	err := loadMaze("resources/maze01.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}
	for {
		printMaze()
		break
	}
}

func loadMaze(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		maze = append(maze, line)
	}
	return nil
}

func printMaze() {
	for _, line := range maze {
		fmt.Println(line)
	}
}
