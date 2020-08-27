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
		if rand.Intn(8) == 1 {
			randomGhostMovement(ghost)
		} else {
			smartGhostMovement(ghost)
		}
	}
	handleGhostContact()
}

func randomGhostMovement(ghostToMove *ghost) {
	direction := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "LEFT",
		3: "RIGHT",
	}
	makeMove(ghostToMove.position.row, ghostToMove.position.col, move[direction])
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

func updateGhostsStatus(ghosts []*ghost, status GhostStatus) {
	ghostStatusMutex.Lock()
	defer ghostStatusMutex.Unlock()
	for _, ghost := range ghosts {
		ghost.status = status
	}
}

type node struct {
	row, col int
	previous *node
}

func smartGhostMovement(ghost *ghost) {
	startingLocation := node{ghost.position.row, ghost.position.col, nil}
	destinationLocation := node{Player.row, Player.col, nil}
	currentLocation := node{ghost.position.row, ghost.position.col, nil}
	var visited []node
	visited = append(visited, startingLocation)
	index := 0

	for currentLocation.row != destinationLocation.row || currentLocation.col != destinationLocation.col {
		currentLocation = visited[index]
		adjacentNodes := findAdjacentNodes(currentLocation)
		for _, adjacentNode := range adjacentNodes {
			if !contains(visited, adjacentNode) {
				visited = append(visited, adjacentNode)
			}
		}
		index++
	}

	for *currentLocation.previous != startingLocation {
		currentLocation = *currentLocation.previous
	}
	ghost.position.row, ghost.position.col = currentLocation.row, currentLocation.col
}

func findAdjacentNodes(fromNode node) []node {
	var adjacentLocations []node
	grid := maze.toGrid()

	possibleMoves := []string{"UP", "DOWN", "LEFT", "RIGHT"}
	for _, move := range possibleMoves {
		switch move {
		case "UP":
			if validLocation(grid, fromNode.row-1, fromNode.col) {
				adjacentLocations = append(adjacentLocations, node{fromNode.row - 1, fromNode.col, &fromNode})
			}
		case "DOWN":
			if validLocation(grid, fromNode.row+1, fromNode.col) {
				adjacentLocations = append(adjacentLocations, node{fromNode.row + 1, fromNode.col, &fromNode})
			}
		case "LEFT":
			if validLocation(grid, fromNode.row, fromNode.col-1) {
				adjacentLocations = append(adjacentLocations, node{fromNode.row, fromNode.col - 1, &fromNode})
			}
		case "RIGHT":
			if validLocation(grid, fromNode.row, fromNode.col+1) {
				adjacentLocations = append(adjacentLocations, node{fromNode.row, fromNode.col + 1, &fromNode})
			}
		}
	}
	return adjacentLocations
}

func validLocation(grid [][]rune, row, col int) bool {
	if row < 0 || col < 0 || row >= len(grid) || col >= len(grid[0]) || grid[row][col] == '#' {
		return false
	} else {
		return true
	}
}

func contains(visitedNodes []node, nodeToFind node) bool {
	for _, n := range visitedNodes {
		if n.col == nodeToFind.col && n.row == nodeToFind.row {
			return true
		}
	}
	return false
}
