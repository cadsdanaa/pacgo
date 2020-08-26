package main

import "math/rand"

type GhostStatus string

const (
	Normal GhostStatus = "Normal"
	Blue   GhostStatus = "Blue"
)

type ghost struct {
	position sprite
	status   GhostStatus
}

func MoveGhosts() {
	handleGhostContact()
	for _, ghost := range ghosts {
		direction := ghostMovement()
		ghost.position.row, ghost.position.col, _ = makeMove(ghost.position.row, ghost.position.col, direction)
	}
	handleGhostContact()
}

func handleGhostContact() {
	for _, ghost := range ghosts {
		if Player.row == ghost.position.row && Player.col == ghost.position.col {
			if ghost.status == Normal {
				handleNormalGhostContact()
			} else {
				handleBlueGhostContact(ghost)
			}
		}
	}
}

func (ghost *ghost) ResetPosition() {
	ghost.position.row, ghost.position.col = ghost.position.startRow, ghost.position.startCol
	ghost.status = Normal
}

func ResetGhosts() {
	for _, ghost := range ghosts {
		ghost.ResetPosition()
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

func updateGhostsStatus(ghosts []*ghost, status GhostStatus) {
	ghostStatusMutex.Lock()
	defer ghostStatusMutex.Unlock()
	for _, ghost := range ghosts {
		ghost.status = status
	}
}
