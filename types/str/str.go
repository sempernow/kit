// Package str provides string-related functions
package str

import (
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandAlphaNum is an efficient random alphanumeric generator.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandAlphaNum(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// Reverse a string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// StripQuotes from a string
// https://stackoverflow.com/questions/44222554/how-to-remove-quotes-from-around-a-string-in-golang#44222606
func StripQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

// Unique returns input slice stripped of redundant elements
func Unique(ss []string) []string {
	sort.Strings(ss)
	j := 0
	for i := 1; i < len(ss); i++ {
		if ss[j] == ss[i] {
			continue
		}
		j++
		ss[j] = ss[i]
	}
	return ss[:j+1]
}

func StripDuplSpace(data string) string {
	ws := regexp.MustCompile(`\s+`)
	return ws.ReplaceAllString(data, " ")
}

// IsAlphaNum ...
func IsAlphaNum(s string) bool {
	safe := func(r rune) rune {
		switch {
		case r > unicode.MaxASCII:
			return -1
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		default:
			return -1
		}
	}
	len0 := len(s)
	s = strings.Map(safe, s)
	return len(s) == len0
}

// CleanAlphaNum filters out all characters
// but for alphanum words, each separated by single whitespace.
// Optionally limit return to max characters.
func CleanAlphaNum(s string, max ...int) string {
	safe := func(r rune) rune {
		switch {
		case r > unicode.MaxASCII:
			return -1
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		case rune(' ') == r:
			return r
		default:
			return -1
		}
	}
	s = strings.Map(safe, s)
	//s = strings.Title(strings.ToLower(s))
	nn := strings.Fields(s)
	s = strings.Join(nn[:], " ")

	if len(max) > 0 {
		if len(s) > max[0] {
			s = s[:max[0]]
		}
	}
	return s
}

// UniqueStrings removes duplicates.
func UniqueStrings(ss []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, s := range ss {
		if _, val := keys[s]; !val {
			keys[s] = true
			list = append(list, s)
		}
	}
	return list
} // https://www.golangprograms.com/remove-duplicate-values-from-slice.html
