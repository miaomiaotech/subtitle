package subtitle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sb "github.com/martinlindhe/subtitles"
)

func Load(path string) (*sb.Subtitle, error) {
	inputText, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file failed: %v", err)
	}

	var sub sb.Subtitle
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".srt":
		sub, err = sb.NewFromSRT(string(inputText))
	case ".ssa":
		sub, err = sb.NewFromSSA(string(inputText))
	case ".vtt":
		sub, err = sb.NewFromVTT(string(inputText))
	case ".ccdb":
		sub, err = sb.NewFromCCDBCapture(string(inputText))
	default:
		return nil, fmt.Errorf("unsupported subtitle format %s", ext)
	}

	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func Dump(sub *sb.Subtitle, path string) error {
	return os.WriteFile(path, []byte(sub.AsSRT()), 0644)
}

func PathAddSuffix(path, suffix, newExt string) string {
	ext := filepath.Ext(path)
	return fmt.Sprintf("%s.%s.%s", strings.TrimSuffix(path, ext), suffix, newExt)
}
