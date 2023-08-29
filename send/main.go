package main

import (
	"github.com/linxGnu/gosmpp/data"
	"log"
)

//var messages = []string{
//	"@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞ^{}\\[~]|€ÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüàABCDEFGHIJKLMNOPQA",
//}

func main() {
	msg := "jdhfkasdfhkjasdhfkjadshfkjadshfkjadshfkjadshfkjadshfjkadshfjkasdhfkjadshfakjdshfajksdfhajkdsfhajksdhfjkadshfkjadshfjkadshfjkasdfhaksdjfhdjk1"

	a, b := data.SplitSms(msg)

	data.FindInAsciiNonGsm38()
	log.Println(a)
	log.Println(b.DataCoding())
	//a, b, c, d := data2.SMSParts("Salom123")
	//aa, bb := d.Encode("Salom123")
	//log.Println(a, b, c, d.DataCoding())
	//log.Println(aa, bb)
	//runes := []rune(messages[0])
	//log.Println(utf8.RuneCountInString(messages[0]), len(runes))
}
