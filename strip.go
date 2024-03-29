package mist

import "strings"

func StripComments(code string) string {
	// TODO: Mist doesn't support strings or any other form of
	// escaping ';', so this function is good enough.

	var out strings.Builder

	comment := false
	for _, c := range code {
		if c == ';' {
			comment = true
		} else if c == '\n' {
			comment = false
		}

		if !comment {
			out.WriteRune(c)
		}
	}

	return out.String()
}
