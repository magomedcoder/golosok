package audio

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/magomedcoder/golosok/internal/core"
)

var rhvoiceTestCandidates = []string{"RHVoice-test", "rhvoice-test"}

var rhvoiceTestExe string

func RegisterRHVoice(c *core.Core) {
	c.RegisterTTS("rhvoice", TTSRHVoiceInit, nil, TTSRHVoiceFile)
}

func findRHVoiceTest() (string, error) {
	for _, name := range rhvoiceTestCandidates {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("в PATH нет RHVoice-test (пробовали: %s). Установите пакет rhvoice (Debian/Ubuntu) или добавьте каталог с бинарником в PATH", strings.Join(rhvoiceTestCandidates, ", "))
}

func rhvoiceProfile() string {
	if v := strings.TrimSpace(os.Getenv("GOLOSOK_RHVOICE_PROFILE")); v != "" {
		return v
	}

	return "anna"
}

func TTSRHVoiceInit(*core.Core) error {
	p, err := findRHVoiceTest()
	if err != nil {
		return err
	}

	rhvoiceTestExe = p
	log.Printf("TTS RHVoice: %s, профиль %q (переопределение: GOLOSOK_RHVOICE_PROFILE)", p, rhvoiceProfile())
	return nil
}

func TTSRHVoiceFile(_ *core.Core, text, out string) error {
	exe := rhvoiceTestExe
	if exe == "" {
		var err error
		exe, err = findRHVoiceTest()
		if err != nil {
			return err
		}
	}

	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return fmt.Errorf("rhvoice: каталог для wav: %w", err)
	}

	prof := rhvoiceProfile()

	args := []string{"-p", prof, "-R", "24000", "-o", out}
	cmd := exec.Command(exe, args...)
	cmd.Stdin = bytes.NewReader([]byte(text))

	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(outBytes))
		if msg != "" {
			return fmt.Errorf("rhvoice (%s %v): %w: %s", exe, args, err, msg)
		}

		return fmt.Errorf("rhvoice (%s %v): %w", exe, args, err)
	}

	return nil
}
