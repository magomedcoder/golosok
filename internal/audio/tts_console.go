package audio

import "fmt"

func RegisterConsole(c *Core) {
	c.RegisterTTS("console", TTSConsoleInit, TTSConsoleSay)
}

func TTSConsoleInit(*Core) error {
	fmt.Println("TTS init console")
	return nil
}

func TTSConsoleSay(c *Core, s string) error {
	fmt.Printf("TTS console: %s", s)
	return nil
}
