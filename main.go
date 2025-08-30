package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	mic, err := NewMic()
	if err != nil {
		log.Fatalf("mic: %v", err)
	}
	defer func(mic *Mic) {
		if err := mic.Close(); err != nil {
			log.Printf("mic-сlose: %v", err)
		}
	}(mic)

	buf := make([]byte, 8000)
	for ctx.Err() == nil {

		n, err := mic.Read(buf)
		if err != nil {
			continue
		}

		fmt.Println(buf[:n])
	}
}
