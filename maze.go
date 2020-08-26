package main

import (
	"bufio"
	"fmt"
	"github.com/danicat/simpleansi"
	"math/rand"
	"os"
	"time"
)

type sprite struct {
	row      int
	col      int
	startRow int
	startCol int
}

type GhostStatus string

const (
	Normal GhostStatus = "Normal"
	Blue   GhostStatus = "Blue"
)

type ghost struct {
	position sprite
	status   GhostStatus
}

var maze []string
var ghosts []*ghost
var score int
var Player sprite
var Lives = 3
var Dots int

func LoadMaze(file string) error {
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
	loadSprites()
	return nil
}

func loadSprites() {
	for row, line := range maze {
		for col, character := range line {
			switch character {
			case 'P':
				Player = sprite{row, col, row, col}
			case 'G':
				ghosts = append(ghosts, &ghost{sprite{row, col, row, col}, Normal})
			case '.', 'X':
				Dots++
			}
		}
	}
}

func PrintMaze() {
	simpleansi.ClearScreen()
	for _, line := range maze {
		for _, character := range line {
			switch character {
			case '#':
				fmt.Print(simpleansi.WithBlueBackground(Cfg.Wall))
			case '.':
				fmt.Print(Cfg.Dot)
			case 'X':
				fmt.Print(Cfg.Pill)
			default:
				fmt.Print(Cfg.Space)
			}
		}
		fmt.Println()
	}
	moveCursor(Player.row, Player.col)
	fmt.Print(Cfg.Player)
	for _, ghost := range ghosts {
		moveCursor(ghost.position.row, ghost.position.col)
		if ghost.status == Normal {
			fmt.Print(Cfg.Ghost)
		} else {
			fmt.Print(Cfg.BlueGhost)
		}
	}

	moveCursor(len(maze)+1, 0)
	fmt.Println("Score: ", score, "\tLives: ", getLivesIcon(), "\tPower Pill Time: ", pillTimer.timeLeft().Truncate(time.Second/10))
}

func MoveGhosts() {
	for _, ghost := range ghosts {
		direction := ghostMovement()
		ghost.position.row, ghost.position.col, _ = makeMove(ghost.position.row, ghost.position.col, direction)
	}
}

func ghostMovement() string {
	direction := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "LEFT",
		3: "RIGHT",
	}
	return move[direction]
}

func moveCursor(row, col int) {
	if Cfg.UseEmoji {
		simpleansi.MoveCursor(row, col*2)
	} else {
		simpleansi.MoveCursor(row, col)
	}
}

func getLivesIcon() string {
	var currentLives = ""
	for i := 0; i < Lives; i++ {
		currentLives += Cfg.Player
	}
	return currentLives
}
