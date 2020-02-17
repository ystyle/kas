package character

import (
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var decoder *encoding.Decoder

func init() {

	decoder = simplifiedchinese.GBK.NewDecoder()
}

func preNUm(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}
func IsUtf8(data []byte) bool {
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			return false
		}
	}
	return true
}

func ToUTF8(str string) string {
	buff, err := decoder.Bytes([]byte(str))
	if err != nil {
		return str
	}
	return string(buff)
}
