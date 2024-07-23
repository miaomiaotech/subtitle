package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/miaomiaotech/subtitle"
	"github.com/xwjdsh/fy"
)

const (
	timeFmt = "15:04:05.000"
)

func main() {
	input := flag.String("i", "/dev/stdin", "the input subtitle file path")
	fromLang := flag.String("from", fy.English, "source language")
	toLang := flag.String("to", fy.Chinese, "target language")
	both := flag.Bool("both", false, "output both language")
	flag.Parse()

	sub, err := subtitle.Translate(*input, *fromLang, *toLang, *both)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var outPath string
	if *both {
		outPath = subtitle.PathAddSuffix(outPath, fmt.Sprintf("%s-%s", *fromLang, *toLang), "srt")
	} else {
		outPath = subtitle.PathAddSuffix(*input, *toLang, "srt")
	}

	if err := subtitle.Dump(sub, outPath); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	log.Printf("translate success, saved to %s", outPath)
}
