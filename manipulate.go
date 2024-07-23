package subtitle

import (
	sb "github.com/martinlindhe/subtitles"
)

func Merge(subA, subB *sb.Subtitle) *sb.Subtitle {
	maxLen := max(len(subA.Captions), len(subB.Captions))
	newSub := sb.Subtitle{Captions: make([]sb.Caption, maxLen)}

	for index := range maxLen {
		var textArrA, textArrB []string
		var meta *sb.Caption
		if index < len(subA.Captions) {
			meta = &subA.Captions[index]
			textArrA = meta.Text
		}
		if index < len(subB.Captions) {
			if meta != nil {
				meta = &subB.Captions[index]
			}
			textArrB = subB.Captions[index].Text
		}

		if len(textArrA) == 0 && len(textArrB) == 0 {
			break
		}

		var lines []string
		lines = append(lines, textArrA...)
		lines = append(lines, textArrB...)

		newCaption := sb.Caption{Seq: meta.Seq, Start: meta.Start, End: meta.End, Text: lines}
		newSub.Captions[index] = newCaption
	}
	return &newSub
}
