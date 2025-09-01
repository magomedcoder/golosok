package greetings

import (
	"fmt"
	"github.com/magomedcoder/golosok/internal/core"
	"math/rand"
	"time"
)

func Register(c *core.Core) {
	c.RegisterCommand("привет", func(c *core.Core, phrase string) {
		sayRand(c, "И тебе привет")
	})

	c.RegisterCommand("дата", func(c *core.Core, phrase string) {
		now := time.Now()
		wd := []string{"понедельник", "вторник", "среда", "четверг", "пятница", "суббота", "воскресенье"}[int(now.Weekday()+6)%7]
		c.Say("сегодня " + wd + ", " + fmtDateRu(now))
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
	days := []string{"первое", "второе", "третье", "четвёртое", "пятое", "шестое", "седьмое", "восьмое", "девятое", "десятое", "одиннадцатое", "двенадцатое", "тринадцатое", "четырнадцатое", "пятнадцатое", "шестнадцатое", "семнадцатое", "восемнадцатое", "девятнадцатое", "двадцатое", "двадцать первое", "двадцать второе", "двадцать третье", "двадцать четвёртое", "двадцать пятое", "двадцать шестое", "двадцать седьмое", "двадцать восьмое", "двадцать девятое", "тридцатое", "тридцать первое"}
	months := []string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}

	return fmt.Sprintf("%s %s", days[t.Day()-1], months[int(t.Month())-1])
}
