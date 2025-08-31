package greetings

import (
	"github.com/magomedcoder/golosok/internal"
	"math/rand"
)

func Register(c *internal.Core) {
	c.RegisterCommand("привет", func(c *internal.Core, phrase string) {
		sayRand(c, "И тебе привет")
	})
}

func sayRand(c *internal.Core, variants ...string) {
	c.Say(variants[rand.Intn(len(variants))])
}
