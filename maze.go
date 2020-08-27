package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/danicat/simpleansi"
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

type sprite struct {
	row      int
	col      int
	startRow int
	startCol int
}

type Maze []string

var maze Maze
var ghosts []*ghost
var score int
var Player sprite
var Lives = 3
var Dots int
var Cfg Config

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

func RemoveDot(row, col int) {
	maze[row] = maze[row][0:col] + " " + maze[row][col+1:]
}

func (maze Maze) toGrid() [][]rune {
	grid := make([][]rune, len(maze))
	for i, _ := range grid {
		grid[i] = []rune(maze[i])
	}
	return grid
}

func makeMove(oldRow, oldCol int, direction string) (newRow, newCol int, validMove bool) {
	newRow, newCol = oldRow, oldCol

	switch direction {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(maze) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(maze) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(maze[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(maze[0]) - 1
		}
	}

	validMove = true
	if maze[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
		validMove = false
	}

	return newRow, newCol, validMove
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

func LoadMazeConfig(file string) error {
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
