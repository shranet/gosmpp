package data

import "unicode"

// FindEncoding returns suitable encoding for a string.
// If string is ascii, then GSM7Bit. If not, then UCS2.
func FindEncoding(s string) (enc Encoding) {
	if IsASCII(s) {
		enc = GSM7BIT
	} else {
		enc = UCS2
	}
	return
}

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
