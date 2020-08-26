package main

import (
	"fmt"
	"time"
)

var previousDirection string

func MovePlayer(direction string) {
	var validMove bool
	Player.row, Player.col, validMove = makeMove(Player.row, Player.col, direction)
	if validMove {
		previousDirection = direction
	}

	switch maze[Player.row][Player.col] {
	case '.':
		Dots--
		score++
		RemoveDot(Player.row, Player.col)
		go playPillSound()
	case 'X':
		Dots--
		score += 10
		RemoveDot(Player.row, Player.col)
		go playPillSound()
		go processPill()
	}
}

func handleNormalGhostContact() {
	Lives -= 1
	go playDeathSound()
	if Lives != 0 {
		moveCursor(Player.row, Player.col)
		fmt.Print(Cfg.Death)
		time.Sleep(2000 * time.Millisecond)
		moveCursor(len(maze)+2, 0)
		Player.row, Player.col = Player.startRow, Player.startCol
		ResetGhosts()
	}
}

func handleBlueGhostContact(ghost *ghost) {
	ghost.ResetPosition()
	score += 100
}
