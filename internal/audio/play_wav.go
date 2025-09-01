package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/magomedcoder/golosok/internal/core"
	"log"
	"os"
	"time"

	"github.com/go-audio/wav"
	"github.com/hajimehoshi/oto/v2"
)

func RegisterWAVPlayer(c *core.Core) {
	c.RegisterPlayWav("oto", WAVPlayerInit, WAVPlayerPlay)
}

var wavPlayer *WAVPlayer

func WAVPlayerInit(*core.Core) error {
	var err error
	wavPlayer, err = NewWAVPlayer()
	return err
}

func WAVPlayerPlay(c *core.Core, path string) error {
	return wavPlayer.Play(path)
}

type WAVPlayer struct {
	ctx          *oto.Context
	sampleRate   int
	channelCount int
}

func NewWAVPlayer() (*WAVPlayer, error) {
	return &WAVPlayer{}, nil
}

func (p *WAVPlayer) ensureContext(sampleRate, channels int) error {
	if p.ctx != nil && p.sampleRate == sampleRate && p.channelCount == channels {
		return nil
	}

	ctx, ready, err := oto.NewContext(sampleRate, channels, oto.FormatSignedInt16LE)
	if err != nil {
		return err
	}

	<-ready

	p.ctx = ctx
	p.sampleRate = sampleRate
	p.channelCount = channels

	return nil
}

func (p *WAVPlayer) Play(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := wav.NewDecoder(f)
	if !dec.IsValidFile() {
		return errors.New("не WAV-файл")
	}

	if dec.BitDepth != 16 {
		return errors.New("поддерживается только 16-битный PCM WAV")
	}

	pcmBuf, err := dec.FullPCMBuffer()
	if err != nil {
		return err
	}

	ints := pcmBuf.AsIntBuffer().Data
	raw := make([]byte, len(ints)*2)
	for i, v := range ints {
		if v > 32767 {
			v = 32767
		} else if v < -32768 {
			v = -32768
		}

		binary.LittleEndian.PutUint16(raw[i*2:], uint16(int16(v)))
	}

	if err := p.ensureContext(int(dec.SampleRate), int(dec.NumChans)); err != nil {
		return err
	}

	player := p.ctx.NewPlayer(bytes.NewReader(raw))
	defer func(player oto.Player) {
		if err := player.Close(); err != nil {
			log.Printf("wav-player-сlose: %v", err)
		}
	}(player)

	player.Play()
	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}
