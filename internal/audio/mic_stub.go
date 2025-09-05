//go:build !portaudio

package audio

type Mic struct{}

func NewMic(sampleRate int) (*Mic, error) {
	return &Mic{}, nil
}

func (m *Mic) SetBlockFunc(f func() bool) {

}

func (m *Mic) Read(dst []byte) (int, error) {
	return 0, nil
}

func (m *Mic) Close() error {
	return nil
}
