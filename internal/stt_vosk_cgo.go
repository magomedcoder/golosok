//go:build cgo_vosk
// +build cgo_vosk

package internal

/*
#cgo CFLAGS:  -I${SRCDIR}/../third_party/vosk/include
#cgo LDFLAGS: -L${SRCDIR}/../third_party/vosk/lib -lvosk -pthread -Wl,-rpath,'$ORIGIN/../third_party/vosk/lib'
#include <vosk_api.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

type VoskSTT struct {
	model *C.VoskModel
	rec   *C.VoskRecognizer
	out   chan string
}

func NewVoskSTT(modelPath string, sampleRate int) *VoskSTT {
	mp := C.CString(modelPath)
	defer C.free(unsafe.Pointer(mp))

	m := C.vosk_model_new(mp)
	if m == nil {
		return nil
	}

	rec := C.vosk_recognizer_new(m, C.float(sampleRate))
	if rec == nil {
		C.vosk_model_free(m)
		return nil
	}

	v := &VoskSTT{
		model: m,
		rec:   rec,
		out:   make(chan string, 8),
	}
	return v
}

func (v *VoskSTT) Accept(pcm []byte) error {
	if len(pcm) == 0 {
		return nil
	}

	res := C.vosk_recognizer_accept_waveform(v.rec, unsafe.Pointer(&pcm[0]), C.int(len(pcm)))
	if int(res) != 0 {
		js := C.vosk_recognizer_result(v.rec)
		if js != nil {
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
		js := C.vosk_recognizer_final_result(v.rec)
		if js != nil {
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
