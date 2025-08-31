package main

import (
	"context"
	"fmt"
	"github.com/magomedcoder/golosok/internal"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sampleRate := 16000

	mic, err := internal.NewMic(sampleRate)
	if err != nil {
		log.Fatalf("mic init: %v", err)
	}
	defer func(mic *internal.Mic) {
		if err := mic.Close(); err != nil {
			log.Printf("mic-сlose: %v", err)
		}
	}(mic)

	stt, err := internal.NewVoskSTT("./models/vosk", sampleRate)
	if err != nil {
		log.Fatalf("vosk init: %v", err)
	}
	defer stt.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case phrase := <-stt.Out():
				if phrase == "" {
					continue
				}

				fmt.Printf("[РАСПОЗНАНО] %s\n", phrase)
			}
		}
	}()

	buf := make([]byte, 8000)
	for ctx.Err() == nil {

		n, err := mic.Read(buf)
		if err != nil {
			continue
		}

		_ = stt.Accept(buf[:n])
	}
}
