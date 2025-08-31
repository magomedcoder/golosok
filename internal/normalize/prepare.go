package normalize

import (
	"github.com/magomedcoder/golosok/internal"
	"regexp"
	"strings"
)

func RegisterPrepare(c *internal.Core) {
	c.RegisterNormalizer("prepare", initPrepare, normPrepare)
}

func initPrepare(*internal.Core) error {
	return nil
}

func normPrepare(c *internal.Core, text string) string {
	if ok, _ := regexp.MatchString(`^[,.?!;:"() «»'ЁА-Яа-яё\d\s%-]+$`, text); ok {
		return text
	}

	text = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(text, " "))

	return text
}
