package utils

import (
	"regexp"
	"strings"
)

func AllNumToText(s string) string {
	if strings.Contains(s, "-") && strings.Count(s, ".") <= 2 {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			return NumToText(parts[0]) + " тире " + NumToText(parts[1])
		}
	}

	return NumToText(s)
}

func NumToText(s string) string {
	neg := false
	if strings.HasPrefix(s, "-") {
		neg = true
		s = s[1:]
	}

	if strings.Contains(s, ".") {
		p := strings.SplitN(s, ".", 2)
		return decimalToText(p[0], p[1], 2)
	}

	w := num2textInt(parseInt(s))
	if neg {
		return "минус " + w
	}

	return w
}

func decimalToText(intPart, frac string, places int) string {
	return num2textInt(parseInt(intPart)) + " " + num2textInt(parseInt(padRight(frac, places)))
}

func padRight(s string, n int) string {
	for len(s) < n {
		s += "0"
	}

	if len(s) > n {
		s = s[:n]
	}

	return s
}

func parseInt(s string) int {
	re := regexp.MustCompile(`\D`)
	s = re.ReplaceAllString(s, "")
	if s == "" {
		return 0
	}

	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}

	return n
}

func num2textInt(n int) string {
	if n == 0 {
		return Units[0]
	}

	var words []string
	order := 0
	for n > 0 {
		chunk := n % 1000
		if chunk != 0 {
			words = append(chunkWords(chunk, order == 1), orderWord(order, chunk))
		}

		n /= 1000
		order++
	}

	var res []string
	for i := len(words) - 1; i >= 0; i-- {
		if words[i] != "" {
			res = append(res, words[i])
		}
	}

	return join(res, " ")
}

func chunkWords(n int, female bool) []string {
	var out []string
	out = append(out, Hundreds[n/100])
	d := (n / 10) % 10
	u := n % 10

	if d == 1 {
		out = append(out, Teens[u])
		return out
	}

	if Tens[d] != "" {
		out = append(out, Tens[d])
	}

	uWord := unit(u, female)
	if uWord != "" {
		out = append(out, uWord)
	}

	return out
}

func unit(u int, female bool) string {
	switch u {
	case 0:
		return ""
	case 1:
		if female {
			return "одна"
		}
		return "один"
	case 2:
		if female {
			return "две"
		}
		return "два"
	default:
		return Units[u]
	}
}

func orderWord(order int, chunk int) string {
	if order == 0 || chunk == 0 {
		return ""
	}

	switch order {
	case 1:
		if chunk%10 == 1 && chunk%100 != 11 {
			return "тысяча"
		}

		if last := chunk % 10; (last >= 2 && last <= 4) && (chunk%100 < 10 || chunk%100 >= 20) {
			return "тысячи"
		}

		return "тысяч"
	case 2:
		if chunk%10 == 1 && chunk%100 != 11 {
			return "миллион"
		}

		if last := chunk % 10; (last >= 2 && last <= 4) && (chunk%100 < 10 || chunk%100 >= 20) {
			return "миллиона"
		}

		return "миллионов"
	case 3:
		if chunk%10 == 1 && chunk%100 != 11 {
			return "миллиард"
		}

		if last := chunk % 10; (last >= 2 && last <= 4) && (chunk%100 < 10 || chunk%100 >= 20) {
			return "миллиарда"
		}

		return "миллиардов"
	}

	return ""
}

func join(a []string, sep string) string {
	out := ""
	for _, s := range a {
		if s == "" {
			continue
		}

		if out != "" {
			out += sep
		}
		out += s
	}

	return out
}
