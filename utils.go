package mist

import "strings"

func padRight32(hex string) string {
	var b strings.Builder
	b.WriteString(hex)
	for i := len(hex); i < 64; i++ {
		b.WriteRune('0')
	}
	return b.String()
}
