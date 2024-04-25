package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	simpleaudiodecoder "github.com/LuzianU/simple-audio-decoder-go"
)

func main() {
	file := "examples/beat.wav"
	targetSampleRate := 192000
	chunkSize := 1024

	pcm, err := simpleaudiodecoder.NewPcmFromFile(file)
	if err != nil {
		fmt.Println("Failed to create Pcm from file")
		return
	}

	clip, err := simpleaudiodecoder.NewAudioClip(pcm, targetSampleRate, chunkSize)
	if err != nil {
		fmt.Println("Failed to create AudioClip from file")
		return
	}

	var resampled [][]float32

	for {
		buffer, isDone, err := clip.ResampleNext()
		if err != nil {
			panic(err)
		}

		if len(resampled) == 0 {
			// initialize resampled buffer
			resampled = make([][]float32, len(*buffer))
			for i := range *buffer {
				resampled[i] = make([]float32, 0)
			}
		}

		for i := range *buffer {
			resampled[i] = append(resampled[i], (*buffer)[i]...)
		}

		if isDone {
			break
		}
	}

	// free resources (important!)
	clip.Free()
	pcm.Free()

	channels := len(resampled)
	frames := len(resampled[0])

	// interleave channels
	interleaved := make([]float32, channels*frames)
	for i := range frames {
		for c := range channels {
			interleaved[i*channels+c] = resampled[c][i]
		}
	}

	// write interleaved data to file
	output, err := os.Create("interleaved.dat")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	err = binary.Write(output, binary.NativeEndian, interleaved)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wrote interleaved data to interleaved.dat")
}
