package audio

import (
	"fmt"
	"github.com/magomedcoder/golosok/internal/core"
)

func RegisterConsole(c *core.Core) {
	c.RegisterTTS("console", TTSConsoleInit, TTSConsoleSay, nil)
}

func TTSConsoleInit(*core.Core) error {
	fmt.Println("TTS init console")
	return nil
}

func TTSConsoleSay(c *core.Core, s string) error {
	fmt.Printf("TTS console: %s", s)
	return nil
}
