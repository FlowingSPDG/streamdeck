package streamdeck

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

// Image encodes img as a PNG data URL suitable for SetImage.
func Image(i image.Image) (string, error) {
	var b bytes.Buffer

	bw := bufio.NewWriter(&b)
	if _, err := bw.WriteString("data:image/png;base64,"); err != nil {
		return "", err
	}

	w := base64.NewEncoder(base64.StdEncoding, bw)
	if err := png.Encode(w, i); err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	if err := bw.Flush(); err != nil {
		return "", err
	}

	return b.String(), nil
}
