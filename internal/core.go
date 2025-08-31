package internal

import (
	"fmt"
	"sync"
	"time"
)

type Core struct {
	timers    [8]int64
	timersEnd [8]func(*Core)
	ctxMu     sync.Mutex

	micBlockedMu sync.RWMutex
	micBlocked   bool

	Commands map[string]interface{}
}

func NewCore() *Core {
	c := &Core{}
	return c
}

func (c *Core) BlockMic() {
	c.micBlockedMu.Lock()
	c.micBlocked = true
	c.micBlockedMu.Unlock()
}

func (c *Core) UnblockMic() {
	c.micBlockedMu.Lock()
	c.micBlocked = false
	c.micBlockedMu.Unlock()
}

func (c *Core) IsMicBlocked() bool {
	c.micBlockedMu.RLock()
	defer c.micBlockedMu.RUnlock()
	return c.micBlocked
}

func (c *Core) UpdateTimers() {
	now := time.Now().Unix()
	for i := range c.timers {
		if c.timers[i] > 0 && now >= c.timers[i] {
			fmt.Printf("End Timer ID=%d at %s\n", i, time.Unix(now, 0).Format("2006-01-02 15:04:05"))
			c.ClearTimer(i, true)
		}
	}
}
func (c *Core) ClearTimer(i int, runEnd bool) {
	if runEnd && c.timersEnd[i] != nil {
		c.timersEnd[i](c)
	}

	c.timers[i] = -1
	c.timersEnd[i] = nil
}

func splitTokens(s string) []string {
	var out []string
	cur := ""
	for _, r := range s {
		if r == ' ' || r == '\t' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
		} else {
			cur += string(r)
		}
	}

	if cur != "" {
		out = append(out, cur)
	}

	return out
}

func joinTokens(t []string) string {
	res := ""
	for i, s := range t {
		if i > 0 {
			res += " "
		}

		res += s
	}
	return res
}

func (c *Core) callFunc(phrase string, fn interface{}) bool {
	switch f := fn.(type) {
	case func(*Core, string):
		f(c, phrase)
		return true
	case func(*Core):
		f(c)
		return true
	case func(string):
		f(phrase)
		return true
	default:
		return false
	}
}

func startsWithAny(phrase string, variants []string) (string, string, bool) {
	for _, v := range variants {
		if len(phrase) >= len(v) && phrase[:len(v)] == v {
			rest := ""
			if len(phrase) > len(v) {
				rest = phrase[len(v):]
			}

			if len(rest) > 0 && rest[0] == ' ' {
				rest = rest[1:]
			}

			return v, rest, true
		}
	}

	return "", "", false
}

func splitVariants(s string) []string {
	var res []string
	cur := ""
	for _, r := range s {
		if r == '|' {
			res = append(res, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}

	if cur != "" {
		res = append(res, cur)
	}

	return res
}

func (c *Core) executeNext(phrase string, ctx interface{}) bool {
	sw, ok := ctx.(map[string]interface{})
	if !ok {
		return false
	}

	fmt.Println("<< executeNext")
	fmt.Println(sw)
	fmt.Println("executeNext >>")

	if fn, found := sw[phrase]; found {
		return c.callFunc(phrase, fn)
	}

	for k, v := range sw {
		fmt.Println("0-0")
		if _, rest, ok := startsWithAny(phrase, splitVariants(k)); ok && rest != "" {
			return c.callFunc(rest, v)
		}
	}

	return false
}

func (c *Core) RunInputStr(s string) bool {
	if s == "" {
		return false
	}

	tokens := splitTokens(s)
	for i, _ := range tokens {

		rest := joinTokens(tokens[i+1:])
		return c.executeNext(rest, c.Commands)

	}

	return c.executeNext(s, c.Commands)
}
