package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type PillTimer struct {
	timer *time.Timer
	end   time.Time
}

var pillTimer PillTimer
var pillMutex sync.Mutex
var ghostStatusMutex sync.RWMutex

func Initialize() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatal("Unable to initialize terminal")
	}
}

func Cleanup() {
	cookedTerminal := exec.Command("stty", "-cbreak", "echo")
	cookedTerminal.Stdin = os.Stdin

	err := cookedTerminal.Run()
	if err != nil {
		log.Fatal("Unable to restore terminal")
	}
}

func ReadInput() (string, error) {
	buffer := make([]byte, 100)
	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}
		}
	}

	return "", nil
}

func MovePlayer(direction string) bool {
	var validMove bool
	Player.row, Player.col, validMove = makeMove(Player.row, Player.col, direction)

	removeDot := func(row, col int) {
		maze[Player.row] = maze[Player.row][0:Player.col] + " " + maze[Player.row][Player.col+1:]
	}

	switch maze[Player.row][Player.col] {
	case '.':
		Dots--
		score++
		removeDot(Player.row, Player.col)
	case 'X':
		Dots--
		score += 10
		removeDot(Player.row, Player.col)
		go processPill()
	}
	return validMove
}

func processPill() {
	pillMutex.Lock()
	updateGhosts(ghosts, Blue)
	pillTime := time.Second * Cfg.PillDurationSeconds
	if pillTimer.timeLeft() > 0 {
		pillTimer.timer.Stop()
		pillTime += pillTimer.timeLeft()
	}
	pillTimer = PillTimer{
		time.NewTimer(pillTime),
		time.Now().Add(pillTime),
	}
	pillMutex.Unlock()
	<-pillTimer.timer.C
	pillMutex.Lock()
	pillTimer.timer.Stop()
	updateGhosts(ghosts, Normal)
	pillMutex.Unlock()
}

func updateGhosts(ghosts []*ghost, status GhostStatus) {
	ghostStatusMutex.Lock()
	defer ghostStatusMutex.Unlock()
	for _, ghost := range ghosts {
		ghost.status = status
	}
}

func (pillTimer *PillTimer) timeLeft() time.Duration {
	remainingTime := pillTimer.end.Sub(time.Now())
	if remainingTime > 0 {
		return remainingTime
	} else {
		return 0
	}
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
