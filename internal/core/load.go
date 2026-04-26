package core

import "fmt"

func (c *Core) RegisterCommand(pattern string, handler any) {
	if c.Commands == nil {
		c.Commands = map[string]any{}
	}

	c.Commands[pattern] = handler
}

func (c *Core) RegisterTTS(id string, init TTSInitFn, say TTSSayFn, toFile TTSToFileFn) {
	c.TTSEngines[id] = [3]any{init, say, toFile}
}

func (c *Core) RegisterNormalizer(id string, init NormalizerInitFn, normalize NormalizeFn) {
	c.Normalizers[id] = [2]any{init, normalize}
}

func (c *Core) RegisterPlayWav(id string, init PlayWAVInitFn, play PlayWAVFn) {
	c.PlayWavs[id] = [2]any{init, play}
}

func (c *Core) DebugListCommands() {
	fmt.Println("[СПИСОК КОМАНД]")
	for k := range c.Commands {
		fmt.Println(" -", k)
	}
}
