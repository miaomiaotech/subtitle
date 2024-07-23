package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/miaomiaotech/subtitle"
)

var chineseKeywords = []string{"zh", "cn", "simp"}

func main() {
	CHINESE_FIRST := os.Getenv("CHINESE_FIRST") == "1"
	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println("Usage: merge-subtitles <subA> <subB>")
		os.Exit(1)
	case 1:
		fmt.Println("missing subB")
		os.Exit(1)
	default:
		if len(args)%2 != 0 {
			fmt.Println("args must be even")
			os.Exit(1)
		}

		for len(args) > 0 {
			x, y := args[0], args[1]
			args = args[2:]

			yl := strings.ToLower(y)
			if CHINESE_FIRST && containsAny(yl, chineseKeywords) {
				x, y = y, x
			}

			a, err := subtitle.Load(x)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			b, err := subtitle.Load(y)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			fmt.Printf("total %d, %d captions\n", len(a.Captions), len(b.Captions))
			c := subtitle.Merge(a, b)

			outPath := subtitle.PathAddSuffix(x, "merged", "srt")
			if err := subtitle.Dump(c, outPath); err != nil {
				fmt.Println(err)
				os.Exit(3)
			}

			fmt.Println(outPath)
		}
	}
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
