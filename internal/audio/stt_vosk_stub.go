//go:build !cgo_vosk

package audio

import (
	"fmt"
)

type VoskSTT struct{}

func NewVoskSTT(modelPath string, sampleRate int, grammarJSON string) (*VoskSTT, error) {
	_ = grammarJSON
	return nil, fmt.Errorf("vosk недоступен")
}

func (v *VoskSTT) Accept(_ []byte) error {
	return fmt.Errorf("vosk недоступен")
}

func (v *VoskSTT) Out() <-chan string {
	ch := make(chan string)
	close(ch)
	return ch
}

func (v *VoskSTT) Close() error {
	return nil
}
