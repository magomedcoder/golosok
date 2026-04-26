//go:build portaudio

package audio

import (
	"errors"

	"github.com/gordonklaus/portaudio"
)

var errMicBufTooSmall = errors.New("буфер микрофона меньше кадра: нужен len >= FrameBytes()")

type Mic struct {
	stream  *portaudio.Stream
	in      []int16
	blockFn func() bool
}

func (m *Mic) FrameBytes() int {
	return len(m.in) * 2
}

func NewMic(sampleRate int) (*Mic, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}

	m := &Mic{}

	m.in = make([]int16, 8000)
	st, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), len(m.in), m.in)
	if err != nil {
		return nil, err
	}

	m.stream = st
	if err := st.Start(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Mic) SetBlockFunc(f func() bool) {
	m.blockFn = f
}

func (m *Mic) Read(dst []byte) (int, error) {
	if m.blockFn != nil && m.blockFn() {
		return 0, nil
	}

	need := len(m.in) * 2
	if len(dst) < need {
		return 0, errMicBufTooSmall
	}

	if err := m.stream.Read(); err != nil {
		return 0, err
	}

	for i, v := range m.in {
		dst[2*i] = byte(v)
		dst[2*i+1] = byte(v >> 8)
	}

	return need, nil
}

func (m *Mic) Close() error {
	if m.stream != nil {
		_ = m.stream.Stop()
		_ = m.stream.Close()
	}

	if err := portaudio.Terminate(); err != nil {
		return err
	}

	return nil
}
