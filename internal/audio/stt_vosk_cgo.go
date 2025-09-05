//go:build cgo_vosk
// +build cgo_vosk

package audio

/*
#cgo CFLAGS:  -I${SRCDIR}/../../build/lib
#cgo LDFLAGS: -L${SRCDIR}/../../build/lib -lvosk -pthread -Wl,-rpath,'$ORIGIN/../../build/lib'
#include <stdlib.h>
#include <vosk_api.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"unsafe"
)

type VoskSTT struct {
	model *C.VoskModel
	rec   *C.VoskRecognizer
	out   chan string
}

func NewVoskSTT(modelPath string, sampleRate int) (*VoskSTT, error) {
	mp := C.CString(modelPath)
	defer C.free(unsafe.Pointer(mp))

	m := C.vosk_model_new(mp)
	if m == nil {
		return nil, errors.New("не удалось создать модель Vosk")
	}

	rec := C.vosk_recognizer_new(m, C.float(sampleRate))
	if rec == nil {
		C.vosk_model_free(m)
		return nil, errors.New("не удалось создать распознаватель Vosk")
	}

	C.vosk_recognizer_set_words(rec, 1)

	v := &VoskSTT{
		model: m,
		rec:   rec,
		out:   make(chan string, 8),
	}
	return v, nil
}

func (v *VoskSTT) Accept(pcm []byte) error {
	if len(pcm) == 0 {
		return nil
	}

	if len(pcm)%2 != 0 {
		return errors.New("неверная длина PCM16LE буфера)")
	}

	res := C.vosk_recognizer_accept_waveform(v.rec, (*C.char)(unsafe.Pointer(&pcm[0])), C.int(len(pcm)))
	if int(res) != 0 {
		if js := C.vosk_recognizer_result(v.rec); js != nil {
			go pushJSON(v.out, C.GoString(js))
		} else if js := C.vosk_recognizer_partial_result(v.rec); js != nil {
			go pushJSON(v.out, C.GoString(js))
		}
	}
	return nil
}

func (v *VoskSTT) Out() <-chan string {
	return v.out
}

func (v *VoskSTT) Close() error {
	if v.rec != nil {
		if js := C.vosk_recognizer_final_result(v.rec); js != nil {
			pushJSON(v.out, C.GoString(js))
		}

		C.vosk_recognizer_free(v.rec)
		v.rec = nil
	}

	if v.model != nil {
		C.vosk_model_free(v.model)
		v.model = nil
	}

	close(v.out)

	return nil
}

func pushJSON(ch chan string, s string) {
	var obj map[string]any
	if json.Unmarshal([]byte(s), &obj) == nil {
		if t, _ := obj["text"].(string); t != "" {
			ch <- t
		}
	}
}
