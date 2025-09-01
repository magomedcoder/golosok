package audio

import (
	"github.com/magomedcoder/golosok/internal/core"
	"log"
)

func RegisterConsole(c *core.Core) {
	c.RegisterTTS("console", TTSConsoleInit, TTSConsoleSay, nil)
}

func TTSConsoleInit(*core.Core) error {
	log.Println("TTS init console")
	return nil
}

func TTSConsoleSay(c *core.Core, s string) error {
	log.Printf("TTS console: %s", s)
	return nil
}
