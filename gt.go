package gt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Set up Build environment first before starting translations.
type Build struct {
	Origin string                       // the origin env
	Target string                       // the target env
	Index  map[string]map[string]string // the index which contains all keys and strings
}

// T() is a shorthand method for Translate. Ignores errors and strictly returns strings.
func (b *Build) T(key string, a ...interface{}) (t string) {
	t, _ = b.Translate(key, a...)
	return t
}

// Translate() translates a key or string from origin to target.
// Parses augmented sprintf format when additional arguments are given.
func (b *Build) Translate(key string, a ...interface{}) (t string, err error) {

	var o string // origin string

	// Try to find origin string by key or key[:2]
	if b.Index[key][b.Origin] != "" {
		o = b.Index[key][b.Origin]
	} else if b.Index[key][b.Origin[:2]] != "" {
		o = b.Index[key][b.Origin[:2]]
	}

	// If key is not found, try matching strings in origin.
	if o == "" {
		for k, v := range b.Index {
			if key == v[b.Origin] {
				o, key = key, k
				break
			}
		}
	}

	// Try to find target string by key or key[:2]
	if b.Index[key][b.Target] != "" {
		t = b.Index[key][b.Target]
	} else if b.Index[key][b.Target[:2]] != "" {
		t = b.Index[key][b.Target[:2]]
	}

	if o == "" || t == "" {
		return t, errors.New("Couldn't find origin or target string.")
	}

	// When no additional arguments are given, there's nothing left to do.
	if len(a) == 0 {
		return t, err
	}

	// Find verbs in both strings.
	oVerbs, tVerbs := findVerbs(o), findVerbs(t)

	if len(oVerbs) < len(a) || len(tVerbs) < len(a) {
		return t, errors.New("Arguments count is higher than verbs count.")
	}

	if len(oVerbs) == len(tVerbs) {
		return t, errors.New("Verbs count in origin and target string do not match")
	}

	// Swap arguments positions and clean up tags.
	r, _ := regexp.Compile(`(#[\w0-9-_]+)`)
	for i, v := range tVerbs {
		for j, v2 := range oVerbs {
			if v == v2 {
				a[j], a[i] = a[i], a[j]
				c := r.ReplaceAllLiteralString(v, "")
				t = strings.Replace(t, v, c, -1)
				break
			}
		}
	}

	t = fmt.Sprintf(t, a...)

	return t, err
}

// findVerbs() finds all occurences of printf verbs with optional tags
func findVerbs(s string) (v []string) {
	r, _ := regexp.Compile(`%(?:\d+\$)?[+-]?(?:[ 0]|'.{1})?-?\d*(?:\.\d+)?#?[bcdeEfFgGopqstTuUvxX%]?[#[\w0-9-_]+]?`)
	m := r.FindAllStringSubmatch(s, -1)
	if len(m[0]) > 0 {
		v = m[0]
	}
	return v
}
