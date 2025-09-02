package normalize

import (
	"github.com/magomedcoder/golosok/internal/core"
	"github.com/magomedcoder/golosok/internal/utils"
	"regexp"
	"strings"
)

func RegisterNumbers(c *core.Core) {
	c.RegisterNormalizer("numbers", initNumbers, normNumbers)
}

func initNumbers(*core.Core) error {
	return nil
}

func normNumbers(c *core.Core, text string) string {
	reDia := regexp.MustCompile(`\d*\.\d+-\d*\.\d+`)
	text = reDia.ReplaceAllStringFunc(text, func(x string) string {
		return utils.AllNumToText(x)
	})

	re := regexp.MustCompile(`-?\d+(?:\.\d+)?`)
	text = re.ReplaceAllStringFunc(text, func(x string) string {
		return utils.AllNumToText(x)
	})

	text = strings.ReplaceAll(text, "%", " процентов")

	return text
}
