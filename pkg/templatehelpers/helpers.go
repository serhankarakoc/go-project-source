package templatehelpers

import (
	"net/url"
	"text/template"
	"time"
)

func TemplateHelpers() template.FuncMap {
	fm := template.FuncMap{
		"CurrentYear": func() int { return time.Now().Year() },
		"Add":         func(a, b int) int { return a + b },
		"Subtract":    func(a, b int) int { return a - b },
		"Mul":         func(a, b int) int { return a * b },
		"Max": func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		"Min": func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		"Iterate": func(start, end int) []int {
			count := end - start + 1
			if count <= 0 {
				return []int{}
			}
			items := make([]int, count)
			for i := 0; i < count; i++ {
				items[i] = start + i
			}
			return items
		},
		"urlquery": func(s string) string { return url.QueryEscape(s) },
		"dict": func(values ...interface{}) map[string]interface{} {
			dict := make(map[string]interface{})
			if len(values)%2 != 0 {
				return dict
			}
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					continue
				}
				dict[key] = values[i+1]
			}
			return dict
		},

		"FormatTime": func(t time.Time, layout string) string {
			if t.IsZero() {
				return ""
			}
			return t.Format(layout)
		},

		"FormatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("02.01.2006")
		},

		"FormatDateTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("02.01.2006 15:04")
		},

		"hasPrefix": func(s, prefix string) bool {
			return len(s) >= len(prefix) && s[:len(prefix)] == prefix
		},

		"EqUintPtr": func(a uint, b *uint) bool {
			if b == nil {
				return false
			}
			return a == *b
		},
	}
	return fm
}
