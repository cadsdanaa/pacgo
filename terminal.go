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

func ReadUserInput() chan string {
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := translateInput()
			if err != nil {
				ch <- "ESC"
				log.Fatal("Failed to read input")
			}
			ch <- input
		}
	}(input)
	return input
}

func translateInput() (string, error) {
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

func processPill() {
	pillMutex.Lock()
	updateGhostsStatus(ghosts, Blue)
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
	updateGhostsStatus(ghosts, Normal)
	pillMutex.Unlock()
}

func (pillTimer *PillTimer) timeLeft() time.Duration {
	remainingTime := pillTimer.end.Sub(time.Now())
	if remainingTime > 0 {
		return remainingTime
	} else {
		return 0
	}
}
