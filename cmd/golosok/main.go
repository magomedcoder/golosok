package main

import (
	"context"
	"flag"
	"github.com/magomedcoder/golosok/internal/audio"
	"github.com/magomedcoder/golosok/internal/commands/greetings"
	"github.com/magomedcoder/golosok/internal/core"
	"github.com/magomedcoder/golosok/internal/normalize"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type STT interface {
	Accept([]byte) error

	Out() <-chan string

	Close() error
}

func main() {
	var sttTest int
	flag.IntVar(&sttTest, "stt-test", 0, "STT test")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	sampleRate := 16000

	c := core.NewCore()

	audio.RegisterConsole(c)
	audio.RegisterRHVoice(c)
	audio.RegisterWAVPlayer(c)

	normalize.RegisterPrepare(c)
	normalize.RegisterNumbers(c)

	greetings.Register(c)

	if err := c.SetupAssistantVoice(); err != nil {
		log.Printf("ошибка при настройке: %v", err)
	}

	var stt STT
	var err error

	if sttTest == 1 {
		lines := []string{"привет", "дата", "команды"}
		stt = audio.NewFakeSTT(lines, 500*time.Millisecond)
	} else {
		stt, err = audio.NewVoskSTT("./models/vosk", sampleRate)
		if err != nil {
			log.Fatalf("vosk init: %v", err)
		}
	}
	defer func() {
		if err := stt.Close(); err != nil {
			log.Printf("stt-close: %v", err)
		}
	}()

	var mic *audio.Mic
	if sttTest != 1 {
		mic, err = audio.NewMic(sampleRate)
		if err != nil {
			log.Fatalf("mic-init: %v", err)
		}
		defer func(mic *audio.Mic) {
			if err := mic.Close(); err != nil {
				log.Printf("mic-сlose: %v", err)
			}
		}(mic)
		mic.SetBlockFunc(c.IsMicBlocked)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case phrase := <-stt.Out():
				if phrase == "" {
					continue
				}

				log.Printf("[РАСПОЗНАНО] %s\n", phrase)
				c.BlockMic()
				c.RunInputStr(phrase)
				c.UnblockMic()
				c.UpdateTimers()
			}
		}
	}()

	if sttTest != 1 {
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

			_ = stt.Accept(buf[:n])
		}
	} else {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.UpdateTimers()
			}
		}
	}
}
