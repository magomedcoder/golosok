package greetings

import (
	"fmt"
	"github.com/magomedcoder/golosok/internal/core"
	"github.com/magomedcoder/golosok/internal/utils"
	"math/rand"
	"strconv"
	"time"
)

func Register(c *core.Core) {
	c.RegisterCommand("привет", func(c *core.Core, phrase string) {
		sayRand(c, "И тебе привет")
	})

	c.RegisterCommand("дата", func(c *core.Core, phrase string) {
		now := time.Now()
		wd := utils.Weeks[int(now.Weekday()+6)%7]
		c.Say("сегодня " + wd + ", " + fmtDateRu(now))
	})

	c.RegisterCommand("время", func(c *core.Core, phrase string) {
		now := time.Now()
		h := now.Hour()
		m := now.Minute()
		if m > 0 {
			c.Say(fmt.Sprintf("Сейчас %s, %s", utils.NumToText(strconv.Itoa(h)), utils.NumToText(strconv.Itoa(m))))
		} else {
			c.Say("Сейчас ровно " + utils.NumToText(strconv.Itoa(h)))
		}
	})

	c.RegisterCommand("команды", func(c *core.Core, phrase string) {
		c.DebugListCommands()
		c.Say("Команды распечатаны в консоль")
	})
}

func sayRand(c *core.Core, variants ...string) {
	c.Say(variants[rand.Intn(len(variants))])
}

func fmtDateRu(t time.Time) string {
	return fmt.Sprintf("%s %s", utils.Days[t.Day()-1], utils.Months[int(t.Month())-1])
}
