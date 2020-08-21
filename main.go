package main

import (
	"bufio"
	"fmt"
	"github.com/danicat/simpleansi"
	"log"
	"os"
	"os/exec"
)

var maze []string

func main() {
	initialize()
	defer cleanup()

	err := loadMaze("resources/maze01.txt")
	if err != nil {
		log.Fatal("Failed to load maze")
		return
	}
	for {
		printMaze()

		input, err := readInput()
		if err != nil {
			log.Fatal("Failed to read input")
			break
		}

		if input == "ESC" {
			break
		}
	}
}

func initialize() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatal("Unable to initialize terminal")
	}
}

func cleanup() {
	cookedTerminal := exec.Command("stty", "-cbreak", "echo")
	cookedTerminal.Stdin = os.Stdin

	err := cookedTerminal.Run()
	if err != nil {
		log.Fatal("Unable to restore terminal")
	}
}

func readInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	}

	return "", nil
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
	simpleansi.ClearScreen()
	for _, line := range maze {
		fmt.Println(line)
	}
}
