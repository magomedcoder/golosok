package audio

import (
	"bytes"
	"fmt"
	"github.com/magomedcoder/golosok/internal/core"
	"os/exec"
)

func RegisterRHVoice(c *core.Core) {
	c.RegisterTTS("rhvoice", TTSRHVoiceInit, nil, TTSRHVoiceFile)
}

func TTSRHVoiceInit(*core.Core) error {
	fmt.Println("TTS init RHVoice")
	return nil
}

func TTSRHVoiceFile(c *core.Core, text, out string) error {
	fmt.Println("TTS init RHVoice")
	cmd := exec.Command("RHVoice-test", "-p", "anna", "-o", out)
	cmd.Stdin = bytes.NewReader([]byte(text))
	return cmd.Run()
}
