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
		return text
	}

	text = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(text, " "))

	return text
}
