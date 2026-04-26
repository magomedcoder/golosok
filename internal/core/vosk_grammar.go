package core

import (
	"encoding/json"
	"sort"
	"strings"
)

func (c *Core) STTGrammarJSON() string {
	if len(c.Commands) == 0 {
		return ""
	}

	seen := make(map[string]struct{})
	add := func(s string) {
		s = strings.TrimSpace(s)
		if s == "" {
			return
		}

		seen[s] = struct{}{}
	}

	for key := range c.Commands {
		for _, part := range splitVariants(key) {
			add(part)
			for _, w := range c.VoiceNames {
				w = strings.TrimSpace(w)
				if w == "" {
					continue
				}
				add(w + " " + part)
			}
		}
	}

	out := make([]string, 0, len(seen)+1)
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	out = append(out, "[unk]")

	b, err := json.Marshal(out)
	if err != nil {
		return ""
	}

	return string(b)
}
