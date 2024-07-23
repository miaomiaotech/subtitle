package main

import (
	"fmt"
	"os"

	"github.com/miaomiaotech/subtitle"
)

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println("Usage: merge-subtitles <subA> <subB>")
		os.Exit(1)
	case 1:
		fmt.Println("missing subB")
		os.Exit(1)
	case 2:
		a, err := subtitle.Load(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		b, err := subtitle.Load(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		fmt.Printf("total %d, %d captions\n", len(a.Captions), len(b.Captions))
		c := subtitle.Merge(a, b)

		outPath := subtitle.PathAddSuffix(args[0], "merged", "srt")
		if err := subtitle.Dump(c, outPath); err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

		fmt.Println(outPath)
	default:
		fmt.Println("too many arguments")
		os.Exit(1)
	}
}
