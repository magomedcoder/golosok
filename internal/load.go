package internal

func (c *Core) RegisterCommand(pattern string, handler interface{}) {
	if c.Commands == nil {
		c.Commands = map[string]interface{}{}
	}

	c.Commands[pattern] = handler
}
