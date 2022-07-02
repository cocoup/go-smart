package common

import (
	"path/filepath"
	"strings"
)

type Pattern map[string]struct{}

func (p Pattern) Match(s string) bool {
	for v := range p {
		match, err := filepath.Match(strings.TrimSpace(v), s)
		if err != nil {
			continue
		}
		if match {
			return true
		}
	}
	return false
}

func (p Pattern) list() []string {
	var ret []string
	for v := range p {
		ret = append(ret, v)
	}
	return ret
}

func ParseTableList(tableValue string) Pattern {
	tablePattern := make(Pattern)

	fields := strings.FieldsFunc(tableValue, func(r rune) bool {
		return r == ','
	})
	for _, f := range fields {
		tablePattern[f] = struct{}{}
	}

	return tablePattern
}
