package main

import (
	"github.com/gordonklaus/portaudio"
)

type Mic struct {
	stream *portaudio.Stream
}

func NewMic() (*Mic, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, err
	}

	m := &Mic{}
	sampleRate := 16000
	in := make([]int16, 8000)
	st, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), len(in), in)
	if err != nil {
		return nil, err
	}

	m.stream = st
	if err := st.Start(); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Mic) Read(dst []byte) (int, error) {
	buf := make([]int16, len(dst)/2)
	if err := m.stream.Read(); err != nil {
		return 0, err
	}

	n := copy(dst, int16SliceToBytes(buf))

	return n, nil
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

func int16SliceToBytes(s []int16) []byte {
	b := make([]byte, len(s)*2)
	for i, v := range s {
		b[2*i] = byte(v)
		b[2*i+1] = byte(v >> 8)
	}

	return b
}
