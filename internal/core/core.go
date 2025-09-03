package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Core struct {
	VoiceNames      []string
	timers          [8]int64
	timersEnd       [8]func(*Core)
	ctxMu           sync.Mutex
	micBlockedMu    sync.RWMutex
	micBlocked      bool
	Commands        map[string]interface{}
	Normalizers     map[string][2]interface{}
	NormalizationID string
	TTSEngines      map[string][3]interface{}
	TTSEngineID     string
	PlayWavID       string
	PlayWavs        map[string][2]interface{}
}

func NewCore() *Core {
	c := &Core{
		VoiceNames:      []string{"голосок", "голос"},
		Commands:        map[string]interface{}{},
		Normalizers:     map[string][2]interface{}{},
		NormalizationID: "prepare",
		TTSEngines:      map[string][3]interface{}{},
		TTSEngineID:     "rhvoice",
		PlayWavID:       "oto",
		PlayWavs:        map[string][2]interface{}{},
	}
	return c
}

func (c *Core) SetupAssistantVoice() error {
	if v, ok := c.TTSEngines[c.TTSEngineID]; ok {
		if init, _ := v[0].(TTSInitFn); init != nil {
			_ = init(c)
		}
	}

	if v, ok := c.Normalizers[c.NormalizationID]; ok {
		if init, _ := v[0].(NormalizerInitFn); init != nil {
			_ = init(c)
		}
	}

	if v, ok := c.PlayWavs[c.PlayWavID]; ok {
		if init, _ := v[0].(PlayWAVInitFn); init != nil {
			_ = init(c)
		}
	}

	return nil
}

func (c *Core) BlockMic() {
	c.micBlockedMu.Lock()
	c.micBlocked = true
	c.micBlockedMu.Unlock()
}

func (c *Core) UnblockMic() {
	c.micBlockedMu.Lock()
	c.micBlocked = false
	c.micBlockedMu.Unlock()
}

func (c *Core) IsMicBlocked() bool {
	c.micBlockedMu.RLock()
	defer c.micBlockedMu.RUnlock()
	return c.micBlocked
}

func (c *Core) SetTimer(seconds int, end func(*Core)) int {
	now := time.Now().Unix()
	for i := range c.timers {
		if c.timers[i] <= 0 {
			c.timers[i] = now + int64(seconds)
			c.timersEnd[i] = end
			log.Printf("Установлен таймер #%d | длительность: %d сек | завершение: %s\n", i, seconds, time.Unix(c.timers[i], 0).Format("2006-01-02 15:04:05"))
			return i
		}
	}

	return -1
}

func (c *Core) UpdateTimers() {
	now := time.Now().Unix()
	for i := range c.timers {
		if c.timers[i] > 0 && now >= c.timers[i] {
			log.Printf("Таймер ID=%d завершен в %s\n", i, time.Unix(now, 0).Format("2006-01-02 15:04:05"))
			c.ClearTimer(i, true)
		}
	}
}

func (c *Core) ClearTimer(i int, runEnd bool) {
	if runEnd && c.timersEnd[i] != nil {
		c.timersEnd[i](c)
	}

	c.timers[i] = -1
	c.timersEnd[i] = nil
}

func (c *Core) ClearTimers() {
	for i := range c.timers {
		if c.timers[i] >= 0 {
			c.ClearTimer(i, false)
		}
	}
}

func splitTokens(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ' ' || r == '\t' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
		} else {
			cur += string(r)
		}
	}

	if cur != "" {
		out = append(out, cur)
	}

	return out
}

func joinTokens(t []string) string {
	res := ""
	for i, s := range t {
		if i > 0 {
			res += " "
		}

		res += s
	}
	return res
}

func (c *Core) callFunc(phrase string, fn interface{}) bool {
	switch f := fn.(type) {
	case func(*Core, string):
		f(c, phrase)
		return true
	case func(*Core):
		f(c)
		return true
	case func(string):
		f(phrase)
		return true
	default:
		return false
	}
}

func startsWithAny(phrase string, variants []string) (string, string, bool) {
	for _, v := range variants {
		if len(phrase) >= len(v) && phrase[:len(v)] == v {
			rest := ""
			if len(phrase) > len(v) {
				rest = phrase[len(v):]
			}

			if len(rest) > 0 && rest[0] == ' ' {
				rest = rest[1:]
			}

			return v, rest, true
		}
	}

	return "", "", false
}

func splitVariants(s string) []string {
	var res []string
	cur := ""
	for _, r := range s {
		if r == '|' {
			res = append(res, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}

	if cur != "" {
		res = append(res, cur)
	}

	return res
}

func hasNested(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

func (c *Core) executeNext(phrase string, ctx interface{}) bool {
	sw, ok := ctx.(map[string]interface{})
	if !ok {
		return false
	}

	if fn, found := sw[phrase]; found {
		return c.callFunc(phrase, fn)
	}

	for k, v := range sw {
		if hasNested(v) {
			if _, rest, ok := startsWithAny(phrase, splitVariants(k)); ok {
				if m, _ := v.(map[string]interface{}); m != nil {
					return c.executeNext(rest, m)
				}
			}
		} else {
			if _, rest, ok := startsWithAny(phrase, splitVariants(k)); ok && rest != "" {
				return c.callFunc(rest, v)
			}
		}
	}

	c.Say("не удалось распознать команду")

	return false
}

func (c *Core) RunInputStr(s string) bool {
	if s == "" {
		return false
	}

	tokens := splitTokens(s)
	for i, t := range tokens {
		for _, vName := range c.VoiceNames {
			if t == vName {
				rest := joinTokens(tokens[i+1:])
				return c.executeNext(rest, c.Commands)
			}
		}
	}

	return c.executeNext(s, c.Commands)
}

func (c *Core) Normalize(text string) string {
	if c.NormalizationID == "none" {
		return text
	}

	if v, ok := c.Normalizers[c.NormalizationID]; ok {
		if fn, _ := v[1].(NormalizeFn); fn != nil {
			return fn(c, text)
		}
	}

	return text
}

func (c *Core) tempFileName() string {
	return filepath.Join("runtime", fmt.Sprintf("core_%d", time.Now().UnixNano()))
}

func (c *Core) sayVia(id string, text string) error {
	if v, ok := c.TTSEngines[id]; ok {
		if say, _ := v[1].(TTSSayFn); say != nil {
			return say(c, text)
		}

		if toFile, _ := v[2].(TTSToFileFn); toFile != nil {
			fName := c.tempFileName() + ".wav"
			if err := toFile(c, text, fName); err != nil {
				return err
			}

			defer os.Remove(fName)

			return c.PlayWav(fName)
		}
	}

	return fmt.Errorf("движок TTS не найден: %s", id)
}

func (c *Core) Say(text string) {
	_ = c.sayVia(c.TTSEngineID, c.Normalize(text))
}

func (c *Core) TTSToFile(text, filename string) error {
	if v, ok := c.TTSEngines[c.TTSEngineID]; ok {
		if toFile, _ := v[2].(TTSToFileFn); toFile != nil {
			return toFile(c, c.Normalize(text), filename)
		}
	}

	return fmt.Errorf("tts в файл не поддерживается")
}

func (c *Core) PlayWav(path string) error {
	if v, ok := c.PlayWavs[c.PlayWavID]; ok {
		if fn, _ := v[1].(PlayWAVFn); fn != nil {
			return fn(c, path)
		}
	}

	return fmt.Errorf("движок play_wav не найден")
}
