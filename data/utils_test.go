package data

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/charmap"
	"log"
	"testing"
)

type testData struct {
	sms            string
	targetTotal    int
	targetSize     int
	targetEncoding byte
}

func TestFindEncoding(t *testing.T) {
	gsm7String := "@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞ^{}\\[~]|€ÆæßÉ!\"#¤%&()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà"
	asciiString := ""
	for i := 32; i < 128; i++ {
		asciiString += string(rune(i))
	}
	latin1String := ""
	for i := 160; i < 256; i++ {
		if e, ok := charmap.ISO8859_1.EncodeRune(rune(i)); ok {
			latin1String += string(e)
		}
	}

	tests := []*testData{
		{gsm7String, 1, 160, GSM7BITCoding},
		{gsm7String + gsm7String, 2, 153, GSM7BITCoding},
		{gsm7String + "`", 2, 63, UCS2Coding},
		//{asciiString, 1, 160, ASCIICoding},
		//{asciiString + asciiString, 2, 153, ASCIICoding},
		//{latin1String, 1, 140, LATIN1Coding},
		//{asciiString + latin1String, 2, 134, LATIN1Coding},
		{asciiString + "р", 2, 63, UCS2Coding},
		{"салом", 1, 70, UCS2Coding},
		{"саломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсалом", 1, 70, UCS2Coding},
		{"саломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломсаломс", 2, 63, UCS2Coding},
	}

	for _, test := range tests {
		defaultEncoding := int16(0)
		if test.targetEncoding == ASCIICoding {
			defaultEncoding = 1
		}
		parts, enc := SplitSms(test.sms, defaultEncoding)
		log.Println(len(parts), enc, test.sms)
		require.Equal(t, len(parts), test.targetTotal, "total")
		//require.Equal(t, parts[0].Bytes, test.targetSize, "size")
		require.Equal(t, enc.DataCoding(), test.targetEncoding, "encoding")
	}

	//require.Equal(t, GSM7BIT, FindEncoding("abc30hb3bk2lopzSD=2-^"))
	//require.Equal(t, UCS2, FindEncoding("Trần Lập và ban nhạc Bức tường huyền thoại"))
	//require.Equal(t, UCS2, FindEncoding("Đừng buồn thế dù ngoài kia vẫn mưa nghiễng rợi tý tỵ"))
}

func TestVonagePartsCount(t *testing.T) {
	type vonadaTestData struct {
		Message    string
		Unicode    bool
		PartsCount int
	}

	vonageTestList := []*vonadaTestData{
		{Message: "Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salo ]",
			PartsCount: 1, Unicode: false},
		{Message: "Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salomlar Salomlar salo ]x",
			PartsCount: 2, Unicode: false},
		{Message: "Salom Приветик смс салои дун Приветик смс салои дун Приветик смс салои дун ыПриветик смс салои дун Приветик смс салои дун Приветик смс салои дун ы Приветик смске",
			PartsCount: 2, Unicode: false},
		{Message: "Salom Приветик смс салои дун Приветик смс салои дун Приветик смс салои дун ыПриветик смс салои дун Приветик смс салои дун Приветик смс салои дун ы Приветик смск",
			PartsCount: 1, Unicode: false},
		{Message: "Salom Приветик смс салои дун Приветик смс салои дун Приветик смс салои дун ыПриветик смс салои дун Приветик смс салои дун Приветик смс салои дун ы Приветик смск",
			PartsCount: 3, Unicode: true},
		{Message: "Salom Приветик смс салои дун Приветик смс салои дун Приветик смс салои дун ыПриветик смс салои дун Приветик смс салои дун Приветик смs",
			PartsCount: 2, Unicode: true},
		{Message: "Salom Приветик смс салои дун Приветик смс салои дун Приветик смс салои",
			PartsCount: 1, Unicode: true},
	}

	for _, td := range vonageTestList {
		require.Equal(t, VonagePartsCount(td.Message, td.Unicode), td.PartsCount, "Parts count")
	}
}
