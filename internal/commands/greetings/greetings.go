package greetings

import (
	"fmt"
	"github.com/magomedcoder/golosok/internal"
)

func Register(c *internal.Core) {
	c.RegisterCommand("привет", func(c *internal.Core, phrase string) {
		fmt.Println("command-greeting")
	})
}
