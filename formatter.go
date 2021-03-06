package gometer

import (
	"bytes"
	"fmt"
	"sort"
)

// Formatter is used to determine a format of metrics representation.
type Formatter interface {
	// Format defines how metrics will be dumped
	// to output destination.
	Format(counters map[string]*Counter) []byte
}

// NewFormatter returns new default formatter.
//
// separator determines how one metric
// will be separated from another.
//
// As separator can be used any symbol: e.g. '\n', ':', '.', ','.
//
// Default format for one metric is: "%v = %v".
// defaultFormatter sorts metrics by value.
func NewFormatter(separator string) Formatter {
	df := &defaultFormatter{
		separator: separator,
	}
	return df
}

type sortedMap struct {
	m map[string]*Counter
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]].Get() < sm.m[sm.s[j]].Get()
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]*Counter) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

type defaultFormatter struct {
	separator string
}

func (f *defaultFormatter) Format(counters map[string]*Counter) []byte {
	var buf bytes.Buffer

	for _, n := range sortedKeys(counters) {
		line := fmt.Sprintf("%v = %v", n, counters[n].Get()) + f.separator
		fmt.Fprintf(&buf, line)
	}

	return buf.Bytes()
}

var _ Formatter = (*defaultFormatter)(nil)
