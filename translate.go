package subtitle

import (
	"context"
	"fmt"
	"log"

	sb "github.com/martinlindhe/subtitles"
	"github.com/xwjdsh/fy"
)

const (
	TimeFmt = "15:04:05.000"
)

func Translate(path string, fromLang, toLang string, both bool) (*sb.Subtitle, error) {
	subtitle, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("load subtitle failed: %v", err)
	}

	newSub := sb.Subtitle{Captions: make([]sb.Caption, len(subtitle.Captions))}
	ctx := context.Background()

	log.Printf("total %d captions, start translating...", len(subtitle.Captions))
	parallel := make(chan bool, 1)
	for i, caption := range subtitle.Captions {
		parallel <- true
		go func() {
			defer func() { <-parallel }()

			log.Printf("\tseq %v, at %v", caption.Seq, caption.Start.Format(TimeFmt))
			var lines []string
			for _, line := range caption.Text {
				if both {
					lines = append(lines, line)
				}
				res := fy.SogouTranslate(ctx, fy.Request{FromLang: fromLang, ToLang: toLang, Text: line})
				if res.Err != nil {
					log.Printf("\t\ttranslate failed, seq %v, at %v, err: %v", caption.Seq, caption.Start.Format(TimeFmt), res.Err)
					continue
				}
				log.Printf("\t\t%v -> %v", line, res.Result)
				lines = append(lines, res.Result)
			}

			newCaption := sb.Caption{Seq: caption.Seq, Start: caption.Start, End: caption.End, Text: lines}
			newSub.Captions[i] = newCaption
		}()
	}
	return &newSub, nil
}
