package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	sub "github.com/martinlindhe/subtitles"
	"github.com/miaomiaotech/sogou"
)

const (
	timeFmt = "15:04:05.000"
)

func main() {
	input := flag.String("i", "/dev/stdin", "the input subtitle file path")
	fromLang := flag.String("from", sogou.English, "source language")
	toLang := flag.String("to", sogou.Chinese, "target language")
	both := flag.Bool("both", false, "output both language")
	flag.Parse()

	inputText, err := os.ReadFile(*input)
	if err != nil {
		log.Printf("read file failed: %v", err)
		os.Exit(1)
	}

	var subtitle sub.Subtitle
	ext := strings.ToLower(filepath.Ext(*input))
	switch ext {
	case ".srt":
		subtitle, err = sub.NewFromSRT(string(inputText))
	case ".ssa":
		subtitle, err = sub.NewFromSSA(string(inputText))
	case ".vtt":
		subtitle, err = sub.NewFromVTT(string(inputText))
	case ".ccdb":
		subtitle, err = sub.NewFromCCDBCapture(string(inputText))
	default:
		fmt.Println("unsupported subtitle format")
		os.Exit(1)
	}

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	xxx := &sogou.Sogou{}
	newSub := sub.Subtitle{Captions: make([]sub.Caption, len(subtitle.Captions))}
	ctx := context.Background()

	log.Printf("total %d captions, start translating...", len(subtitle.Captions))
	parallel := make(chan bool, 3)
	for i, caption := range subtitle.Captions {
		parallel <- true
		go func() {
			defer func() { <-parallel }()

			log.Printf("translating seq %v, at %v", caption.Seq, caption.Start.Format(timeFmt))
			var lines []string
			for _, line := range caption.Text {
				if *both {
					lines = append(lines, line)
				}
				res := xxx.Translate(ctx, sogou.Request{FromLang: *fromLang, ToLang: *toLang, Text: line})
				if res.Err != nil {
					log.Printf("translate failed, seq %v, at %v, err: %v", caption.Seq, caption.Start.Format(timeFmt), res.Err)
					continue
				}
				log.Printf("%v -> %v", line, res.Result)
				lines = append(lines, res.Result)
			}

			newCaption := sub.Caption{Seq: caption.Seq, Start: caption.Start, End: caption.End, Text: lines}
			newSub.Captions[i] = newCaption
		}()
	}

	var newName string
	if *both {
		newName = fmt.Sprintf("%s.%s-%s.srt", strings.TrimSuffix(*input, ext), *fromLang, *toLang)
	} else {
		newName = fmt.Sprintf("%s.%s.srt", strings.TrimSuffix(*input, ext), *toLang)
	}

	if err := os.WriteFile(newName, []byte(newSub.AsSRT()), 0644); err != nil {
		log.Println(err)
		os.Exit(3)
	}
}
