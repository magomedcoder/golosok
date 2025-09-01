package audio

import "time"

type FakeSTT struct {
	out chan string
}

func NewFakeSTT(lines []string, delay time.Duration) *FakeSTT {
	f := &FakeSTT{
		out: make(chan string, len(lines)),
	}
	go func() {
		for _, s := range lines {
			time.Sleep(delay)
			f.out <- s
		}
		close(f.out)
	}()

	return f
}

func (f *FakeSTT) Accept([]byte) error {
	return nil
}

func (f *FakeSTT) Out() <-chan string {
	return f.out
}

func (f *FakeSTT) Close() error {
	return nil
}
