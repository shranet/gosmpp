package data

import (
	"golang.org/x/text/encoding/charmap"
	"log"
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

func FindInAsciiNonGsm38() {
	for i := 32; i < unicode.MaxASCII; i++ {
		char := rune(i)
		if _, ok := forwardLookup[char]; !ok {
			if _, ok := forwardEscape[char]; !ok {
				log.Println(i, string(char))
			}
		}
	}
}

func VonagePartsCount(text string, unicode bool) int {
	if unicode {
		total := utf8.RuneCountInString(text)
		if total <= 70 {
			return 1
		}

		return (total + 66) / 67
	}

	totalSeptets := 0
	for _, char := range text {
		if _, ok := forwardLookup[char]; ok {
			totalSeptets += 1
		} else if _, ok := forwardEscape[char]; ok {
			totalSeptets += 2
		} else {
			//unicode char va vonage 1 byte deb qaraydi
			//odatda ? belgisiga almashtiriladi
			totalSeptets += 1
		}
	}

	if totalSeptets <= 160 {
		return 1
	}

	return (totalSeptets + 152) / 153
}

type SmsPart struct {
	Message string `json:"message"`
	Chars   int    `json:"chars"`
	Bytes   int    `json:"bytes"`
}

func SplitSms(text string, defaultEncoding int16) ([]SmsPart, Encoding) {
	isGsm0338 := true
	isAscii := true
	//isIso88591 := true

	totalSeptets := 0

	for _, char := range text {
		//GSM 0338 emasligiga tekshiryapmiz
		if _, ok := forwardLookup[char]; !ok {
			if _, ok := forwardEscape[char]; !ok {
				isGsm0338 = false
			} else {
				totalSeptets += 2
			}
		} else {
			totalSeptets += 1
		}

		////ASCII emasligiga tekshiryapmiz
		if int(char) > unicode.MaxASCII {
			isAscii = false
		}

		//if _, ok := charmap.ISO8859_1.EncodeRune(char); !ok {
		//	isIso88591 = false
		//}
		//
		//if !isGsm0338 && !isAscii && !isIso88591 {
		//	break
		//}
	}

	//log.Println(isGsm0338, isAscii, isIso88591, totalSeptets)

	if isGsm0338 && defaultEncoding == 0 {
		return splitGsm0338(text, totalSeptets), GSM7BIT
	}

	if isAscii && defaultEncoding == 1 {
		return splitAscii(text), ASCII
	}

	return splitUCS2(text), UCS2
}

func splitGsm0338(text string, totalSeptets int) []SmsPart {
	var result []SmsPart

	if totalSeptets <= 160 {
		return []SmsPart{
			{
				Message: text,
				Bytes:   (totalSeptets*7 + 7) / 8,
				Chars:   utf8.RuneCountInString(text),
			},
		}
	}

	part := SmsPart{Message: "", Chars: 0, Bytes: 0}

	septets := 0
	charSeptet := 0

	for _, char := range text {
		if _, ok := forwardLookup[char]; ok {
			charSeptet = 1
		} else if _, ok := forwardEscape[char]; ok {
			charSeptet = 2
		} else {
			continue
		}

		if septets+charSeptet <= 153 {
			part.Message += string(char)
			part.Chars += 1

			septets += charSeptet
		} else {
			part.Bytes = (septets*7 + 7) / 8
			result = append(result, part)

			part = SmsPart{Message: string(char), Chars: 1, Bytes: 0}
			septets = charSeptet
		}
	}

	if septets > 0 {
		part.Bytes = (septets*7 + 7) / 8
		result = append(result, part)
	}

	return result
}

func splitAscii(text string) []SmsPart {
	var result []SmsPart
	totalSeptets := len([]byte(text))

	if totalSeptets <= 160 {
		return []SmsPart{
			{
				Message: text,
				Bytes:   (totalSeptets*7 + 7) / 8,
				Chars:   utf8.RuneCountInString(text),
			},
		}
	}

	part := SmsPart{Message: "", Chars: 0, Bytes: 0}

	septets := 0
	charSeptet := 1

	for _, char := range text {
		if septets+charSeptet <= 153 {
			part.Message += string(char)
			part.Chars += 1

			septets += charSeptet
		} else {
			part.Bytes = (septets*7 + 7) / 8
			result = append(result, part)

			part = SmsPart{Message: string(char), Chars: 1, Bytes: 0}
			septets = charSeptet
		}
	}

	if septets > 0 {
		part.Bytes = (septets*7 + 7) / 8
		result = append(result, part)
	}

	return result
}

func splitUCS2(text string) []SmsPart {
	var result []SmsPart
	runes := []rune(text)

	if len(runes) <= 70 {
		return []SmsPart{
			{
				Message: text,
				Bytes:   len(runes) * 2,
				Chars:   len(runes),
			},
		}
	}

	for len(runes) > 0 {
		en := 70 - 3
		if en > len(runes) {
			en = len(runes)
		}

		result = append(result, SmsPart{
			Message: string(runes[:en]),
			Bytes:   en * 2,
			Chars:   en,
		})

		runes = runes[en:]
	}

	return result
}
