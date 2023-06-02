package data

import (
	"unicode"
	"unicode/utf8"
)

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

func SMSParts(text string) (int, bool) {
	gsm7octet := GSM7Octet(text)

	isGsm7 := false
	total := 1
	if gsm7octet == -1 {
		smLength := utf8.RuneCountInString(text)
		if smLength > 70 {
			//ceil qilmaslik uchun +
			total = (smLength + (70 - 4)) / (70 - 3) //3 - bu aslida 6 octetning yarmi
		}
	} else {
		isGsm7 = true
		if gsm7octet > 140 {
			//ceil qilmaslik uchun +
			total = (gsm7octet + (140 - 7)) / (140 - 6)
		}
	}

	return total, isGsm7
}
