package normalize

import (
	"github.com/magomedcoder/golosok/internal/core"
	"regexp"
	"strings"
)

func RegisterPrepare(c *core.Core) {
	c.RegisterNormalizer("prepare", initPrepare, normPrepare)
}

func initPrepare(*core.Core) error {
	return nil
}

func normPrepare(c *core.Core, text string) string {
	if ok, _ := regexp.MatchString(`^[,.?!;:"() «»'ЁА-Яа-яё\d\s%-]+$`, text); ok {
		return processNumbers(c, text)
	}

	text = processSymbols(text)
	text = processNumbers(c, text)
	text = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(text, " "))

	return text
}

func processSymbols(s string) string {
	repl := map[rune]string{'%': " процентов "}

	var b strings.Builder
	for _, r := range s {
		if v, ok := repl[r]; ok {
			b.WriteString(v)
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func processNumbers(c *core.Core, s string) string {
	return normNumbers(c, s)
}
