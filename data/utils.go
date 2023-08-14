package data

import (
	"golang.org/x/text/encoding/charmap"
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

func IsISO8859_1(s string) bool {
	for _, r := range s {
		if _, ok := charmap.ISO8859_1.EncodeRune(r); !ok {
			return false
		}
	}
	return true
}

func SMSParts(text string) (int, int, Encoding) {
	gsm7octet := GSM7Octet(text)
	total := 1
	size := -1

	// Agar gsm7 bit bo'lmasa, ASCII ga tekshiramiz
	isAscii := false
	if gsm7octet == -1 {
		isAscii = IsASCII(text)
	}

	//GSM default va ASCII bu 7 bit
	if gsm7octet != -1 || isAscii {
		size = 160 //Bir butun bitta SMS

		smLength := gsm7octet
		if isAscii {
			smLength = len(text)
		}

		if smLength > 160 {
			//Qismlarga ajratilgan SMS. UDH headerni olib tashlaymiz (6 byte)
			//(140 - 6) * 8 / 7 = ~153
			size = 153

			//ceil qilmaslik uchun +
			total = (smLength + (140 - 7)) / (140 - 6)
		}

		var coding Encoding
		if isAscii {
			coding = ASCII
		} else {
			coding = GSM7BIT
		}

		return total, size, coding
	}

	if IsISO8859_1(text) {
		size = 140
		smLength := utf8.RuneCountInString(text)
		if smLength > 140 {
			// 140 - 6 (UDH header size)
			size = 134

			total = (smLength + (140 - 7)) / (140 - 6)
		}

		return total, size, LATIN1
	}

	size = 70
	smLength := utf8.RuneCountInString(text)
	if smLength > 70 {
		size = 63

		//ceil qilmaslik uchun +
		total = (smLength + (70 - 4)) / (70 - 3) //3 - bu aslida 6 octetning yarmi
	}

	return total, size, UCS2
}
