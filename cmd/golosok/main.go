package main

import (
	"context"
	"fmt"
	"github.com/magomedcoder/golosok/internal/audio"
	"github.com/magomedcoder/golosok/internal/commands/greetings"
	"github.com/magomedcoder/golosok/internal/core"
	"github.com/magomedcoder/golosok/internal/normalize"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sampleRate := 16000

	c := core.NewCore()

	audio.RegisterConsole(c)
	audio.RegisterRHVoice(c)

	normalize.RegisterPrepare(c)

	greetings.Register(c)

	if err := c.SetupAssistantVoice(); err != nil {
		log.Printf("ошибка при настройке: %v", err)
	}

	mic, err := audio.NewMic(sampleRate)
	if err != nil {
		log.Fatalf("mic init: %v", err)
	}
	defer func(mic *audio.Mic) {
		if err := mic.Close(); err != nil {
			log.Printf("mic-сlose: %v", err)
		}
	}(mic)

	stt, err := audio.NewVoskSTT("./models/vosk", sampleRate)
	if err != nil {
		log.Fatalf("vosk init: %v", err)
	}
	defer stt.Close()

	mic.SetBlockFunc(c.IsMicBlocked)

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
				c.BlockMic()
				c.RunInputStr(phrase)
				c.UnblockMic()
				c.UpdateTimers()
			}
		}
	}()

	buf := make([]byte, 8000)
	for ctx.Err() == nil {
		c.UpdateTimers()

		if c.IsMicBlocked() {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		n, err := mic.Read(buf)
		if err != nil {
			continue
		}

		//fmt.Println(buf[:n])
		_ = stt.Accept(buf[:n])
	}
}
