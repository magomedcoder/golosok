package core

func (c *Core) RegisterCommand(pattern string, handler interface{}) {
	if c.Commands == nil {
		c.Commands = map[string]interface{}{}
	}

	c.Commands[pattern] = handler
}

func (c *Core) RegisterTTS(id string, init TTSInitFn, say TTSSayFn, toFile TTSToFileFn) {
	c.TTSEngines[id] = [3]interface{}{init, say, toFile}
}

func (c *Core) RegisterNormalizer(id string, init NormalizerInitFn, normalize NormalizeFn) {
	c.Normalizers[id] = [2]interface{}{init, normalize}
}
