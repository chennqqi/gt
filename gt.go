package gt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// var g = &gt.Build{
//  Index: gt.Keys{
//      "homepage-greeting": gt.Strings{
//          "en": "Welcome to %s#title, %s#name!"
//          "es-LA": "¡Bienvenido a %s#title, %s#titlename!",
//          "nl": "Welkom bij %s#title, %s#name!",
//          "tr": "%s#name, %s#title'ya hoşgeldiniz!", // tag notation proves useful in SOV languages where subject comes before object/verb
//          "zh-CN": "欢迎%s#title, %s#name!", 
//      },
//  }, 
//  Origin: "en",
// }
//
// g.Target = "es"
// keyStr := g.T("homepage-greeting", "Github", "John") // outputs: ¡Bienvenido a Github, John!
// 
// g.Target = "tr"
// SOVStr := g.T("Welcome to %s#title, %s#name!", "Yahoo", "Marissa") // outputs: Marissa, Yahoo'ya hoşgeldiniz!
// 
type Strings map[string]string
type Keys map[string]Strings

type Build struct {
	Origin string // the originating env
	Target string // the target env
	Index  Keys   // the index which contains all keys and strings
}

// T is a shorthand method for Translate, ignores errors and strictly returns strings
func (b *Build) T(key string, a ...interface{}) (t string) {
	t, _ = b.Translate(key, a...)
	return t
}

// Translate finds a string for key in target env and transliterates it.
// Will return key if string or target env is not found.
// 
func (b *Build) Translate(key string, a ...interface{}) (t string, err error) {

	var o string // origin string

	o = b.Index[key][b.Origin]

	if o == "" {
		o = b.Index[key][b.Origin[:2]]
	}

	if o == "" { // no key found? try matching strings
		for k, v := range b.Index {
			if key == v[b.Origin] {
				o, key = key, k
				break
			}
		}
	}

	t = b.Index[key][b.Target]

	if t == "" {
		t = b.Index[key][b.Target[:2]]
	}

	if o == "" || t == "" {
		return t, errors.New("Couldn't find key/string")
	}

	if len(a) == 0 { // no arguments? return string
		return t, err
	}

	oVerbs := ParseStr(o)
	tVerbs := ParseStr(t)

	if len(oVerbs) < len(a) || len(tVerbs) < len(a) {
		return t, errors.New("Couldn't find enough verbs to parse args")
	}

	var cleanVerb string

	// time to switch it up!
	if len(oVerbs) == len(tVerbs) { // check if both verbs arrays are the same length
		var r *regexp.Regexp
		r, _ = regexp.Compile(`(#[\w0-9-_]+)`)      // compile regex to match tags
		newArgs := make([]interface{}, len(oVerbs)) // create new args slice

		for ti, dirtyVerb := range tVerbs {
			for oi, ov := range oVerbs {
				if ov == dirtyVerb && a[oi] != nil { // find original argument 
					newArgs[ti] = a[oi]                                  // assign argument to current position
					a[oi] = nil                                          // unset arg
					cleanVerb = r.ReplaceAllLiteralString(dirtyVerb, "") // remove tags
					strings.Replace(t, dirtyVerb, cleanVerb, -1)         // replace dirty verbs with clean verbs
				}
			}
		}

		if len(newArgs) == len(a) {
			t = fmt.Sprintf(t, newArgs...) // replace verbs with args
		} else {
			return t, errors.New("Impossible to assign arguments to string")
		}

	} else {
		return t, errors.New("Number of printf verbs in origin and target string do not match")
	}

	return t, err
}

// ParseStr returns an array of parsed verbs with optional tags
func ParseStr(str string) (verbs []string, err error) {
	r, _ := regexp.Compile(`(%(?:\d+\$)?[+-]?(?:[ 0]|'.{1})?-?\d*(?:\.\d+)?#?[bcdeEfFgGopqstTuUvxX%]?)(#[\w0-9-_]+)?`)
	m := r.FindAllStringSubmatch(str, -1)

	if len(m) > 0 {
		verbs = make([]string, len(m))
		for i, v := range m[0] {
			verbs[i] = v
		}
	}
	return verbs, err
}

// t.Target = "es"

// str := g.T("homepage-greeting")

// t.Target = "nl"

// str = g.T("homepage-greeting")
// Example:
//     fmt.Printf(t.S("Welcome to %s#title, %s!#name!"), "gotr", "Melvin")
//     // outputs: "¡Bienvenido a gotr, Melvin!"
//     
// You can add meta data to the printf verbs by using the tag notation. Not only will translators
// understand more about the context that they're translating in, but it also prevents order issues 
// in certain languages. (The returning string is always correctly parsed).
// 
// Tags are recommended but you can also use regular printf strings. But then the
// order can be wrong in certain SOV languages (Turkish, Japanese):
//
//     fmt.Printf(t.S("Welcome to %s, %s"), "gotr", "Melvin")
//     // outputs: "gotr, Melvin'ya hoşgeldiniz!" in Turkish. 
//     // This roughly translates to: "Welcome to Melvin, gotr!" which is not what you want!
//     
// Tags are normally ended by a single space but you can also end tags with a plus sign:
// "somethingbefore%s#name+somethingafter" 
// 
