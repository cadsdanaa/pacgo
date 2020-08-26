package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

var pillSoundBuffer *beep.Buffer
var deathSoundBuffer *beep.Buffer

func InitSoundBoard() {
	loadPillSound()
	loadDeathSound()
}

func loadPillSound() {
	pillSoundFile, err := os.Open("resources/pill-sound.mp3")
	if err != nil {
		log.Fatal("Failed to load pill sound file")
	}
	streamer, format, err := mp3.Decode(pillSoundFile)
	if err != nil {
		log.Fatal("Failed to decode pill sound file")
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/1000))
	pillSoundBuffer = beep.NewBuffer(format)
	pillSoundBuffer.Append(streamer)
	streamer.Close()
}

func loadDeathSound() {
	deathSoundFile, err := os.Open("resources/pacman-death-sound.mp3")
	if err != nil {
		log.Fatal("Failed to load death sound file")
	}
	streamer, format, err := mp3.Decode(deathSoundFile)
	if err != nil {
		log.Fatal("Failed to decode death sound file")
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/1000))
	deathSoundBuffer = beep.NewBuffer(format)
	deathSoundBuffer.Append(streamer)
	streamer.Close()
}

func playDeathSound() {
	deathSound := deathSoundBuffer.Streamer(0, deathSoundBuffer.Len())
	speaker.Play(deathSound)
}

func playPillSound() {
	pillSound := pillSoundBuffer.Streamer(0, pillSoundBuffer.Len())
	speaker.Play(pillSound)
}
