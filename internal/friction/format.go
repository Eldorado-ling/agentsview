package friction

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.kenn.io/agentsview/internal/stringutil"
)

// truncatedMarker is jilog's suffix.
const truncatedMarker = " … [truncated]"

// TruncateRunes keeps at most n runes and adds no ellipsis.
func TruncateRunes(s string, n int) string {
	return stringutil.TruncateRunes(s, n, "")
}

// TruncateWithMarker keeps at most n runes and appends the marker only
// when text was removed.
func TruncateWithMarker(s string, n int) string {
	return stringutil.TruncateRunes(s, n, truncatedMarker)
}

// WithoutCutTail returns text as a detector of kind k stored it, minus the
// end its length limit may have split: the last token, or a private key
// block left without its END line. A secret split that way no longer
// matches the secret scanner, so text bound off the machine drops it. A
// cut text keeps the truncation marker; other text is returned unchanged.
func WithoutCutTail(k Kind, text string) string {
	head := text
	switch {
	case k == KindError && strings.HasSuffix(text, truncatedMarker):
		head = strings.TrimSuffix(text, truncatedMarker)
	case k == KindWorkaround && utf8.RuneCountInString(text) >= workaroundContextRunes:
	default:
		return text
	}
	if i := strings.LastIndex(head, "-----BEGIN "); i >= 0 && !strings.Contains(head[i:], "-----END ") {
		head = head[:i]
	} else {
		head = head[:max(strings.LastIndexFunc(head, unicode.IsSpace), 0)]
	}
	return strings.TrimRightFunc(head, unicode.IsSpace) + truncatedMarker
}

// PythonRepr approximates Python repr() for a string: single quotes;
// backslash, quote, newline, carriage return, and tab escaped; other
// runes below 0x20 as \xHH; everything else verbatim.
func PythonRepr(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('\'')
	for _, r := range s {
		switch {
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\'':
			b.WriteString(`\'`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20:
			fmt.Fprintf(&b, `\x%02x`, r)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

// SanitizeDisplay replaces Unicode control runes with spaces for one-line display.
func SanitizeDisplay(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsControl(r):
			return ' '
		default:
			return r
		}
	}, s)
}
