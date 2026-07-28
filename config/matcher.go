package config

import (
	"regexp"
	"strings"
)

// Matcher can match exact text or a regular expression.
//
// Regular expressions must be enclosed in slashes (e.g., /pattern/).
type Matcher struct {
	text  string
	regex *regexp.Regexp
}

// Match checks if the given string matches the filter.
func (f *Matcher) Match(s string) bool {
	// Exact match if no regex is defined
	if f.regex == nil {
		return s == f.text
	}

	// Regex match
	return f.regex.MatchString(s)
}

func (f *Matcher) ReplaceAll(s string, replacement string) string {
	if f.regex == nil {
		return replacement
	}

	return f.regex.ReplaceAllString(s, replacement)
}

// UnmarshalText unmarshals the filter from text.
func (f *Matcher) UnmarshalText(text []byte) error {
	var err error

	f.text = string(text)

	// Compile regex if applicable
	if strings.HasPrefix(f.text, "/") && strings.HasSuffix(f.text, "/") && len(f.text) > 2 {
		f.regex, err = regexp.Compile(f.text[1 : len(f.text)-1])
	}

	return err
}
