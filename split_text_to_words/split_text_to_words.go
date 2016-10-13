package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func main() {
	texts := splitTextToWords(Text("中中 ab34c国家bcD12"))
	for _, text := range texts {
		fmt.Println(string(text))
	}
}

type Text []byte

// 将text中的一串字母或数字当做一个Text，其他字当做一个Text输出
func splitTextToWords(text Text) []Text {
	output := make([]Text, 0, len(text)/3)
	current := 0
	inAlphanumeric := true
	alphanumericStart := 0
	for current < len(text) {
		r, size := utf8.DecodeRune(text[current:])
		if size <= 2 && (unicode.IsLetter(r) || unicode.IsNumber(r)) {
			// 当前是拉丁字母或数字（非中日韩文字）
			if !inAlphanumeric {
				// 上一次遍历时不是字母或数字
				alphanumericStart = current
				inAlphanumeric = true
			}
		} else {
			// 中日韩文字
			if inAlphanumeric {
				// 上一次遍历是数字或字母，这次不是了，所以要先将字母或数字串存起来
				inAlphanumeric = false
				if current != 0 {
					output = append(output, text[alphanumericStart:current])
				}
			}
			output = append(output, text[current:current+size]) // 当前指向的rune
		}
		current += size
	}

	// 处理最后一个字元是英文的情况
	if inAlphanumeric {
		if current != 0 {
			output = append(output, text[alphanumericStart:current])
		}
	}

	return output
}
