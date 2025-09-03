package greetings

import (
	"fmt"
	"github.com/magomedcoder/golosok/internal/core"
	"github.com/magomedcoder/golosok/internal/utils"
	"math/rand"
	"strconv"
	"strings"
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

	c.RegisterCommand("таймер|поставь таймер", setTimer)

	c.RegisterCommand("удали таймер|отмени таймер", cancelTimer)

	c.RegisterCommand("удали все таймеры|сбрось все таймеры|отмени все таймеры", func(c *core.Core, phrase string) {
		c.ClearTimers()
		c.Say("Все таймеры остановлены")
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

func setTimer(c *core.Core, phrase string) {
	p := strings.TrimSpace(phrase)
	if p == "" {
		c.SetTimer(
			5*60,
			func(c *core.Core) {
				c.Say("пять минут прошло")
			},
		)

		c.Say("Ставлю таймер на пять минут")
		return
	}

	mins := utils.ParseEndNumber(p, "минут")
	if mins > 0 {
		c.SetTimer(
			mins*60,
			func(c *core.Core) {
				c.Say(utils.NumToText(strconv.Itoa(mins)) + " прошло")
			})

		c.Say("Ставлю таймер на " + utils.NumToText(strconv.Itoa(mins)))
		return
	}

	secs := utils.ParseEndNumber(p, "секунд")
	if secs > 0 {
		c.SetTimer(secs, func(c *core.Core) {
			c.Say(utils.NumToText(strconv.Itoa(secs)) + " прошло")
		})

		c.Say("Ставлю таймер на " + utils.NumToText(strconv.Itoa(secs)))
		return
	}

	c.Say("Что после таймера ?")
}

func cancelTimer(c *core.Core, phrase string) {
	p := strings.TrimSpace(phrase)
	if p == "" {
		c.ClearTimers()
		c.Say("Таймер остановлен")
		return
	}

	n := -1
	for _, f := range strings.Fields(p) {
		if x, err := strconv.Atoi(f); err == nil {
			n = x
			break
		}
	}

	if n >= 0 {
		c.ClearTimer(n, false)
		c.Say("Таймер удалён")
		return
	}

	c.ClearTimers()
	c.Say("Все таймеры остановлены")
}
