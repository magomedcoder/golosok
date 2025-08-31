package greetings

import (
	"github.com/magomedcoder/golosok/internal/core"
	"math/rand"
)

func Register(c *core.Core) {
	c.RegisterCommand("привет", func(c *core.Core, phrase string) {
		sayRand(c, "И тебе привет")
	})
}

func sayRand(c *core.Core, variants ...string) {
	c.Say(variants[rand.Intn(len(variants))])
}
