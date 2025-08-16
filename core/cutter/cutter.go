package cutter

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	_ "image/gif"
	_ "image/jpeg"
)

// CutImageBytes takes an image in raw byte form and cuts it into horizontal segments.
// It returns a slice of image.Image objects, each representing a 1px wide vertical segment.
// The resulting segments can be used for further processing or encoding.
func CutImageBytes(imageData []byte) ([]image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image data: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	sImg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image type does not support SubImage operation")
	}

	segments := make([]image.Image, width)
	for x := 0; x < width; x++ {
		rect := image.Rect(bounds.Min.X+x, bounds.Min.Y, bounds.Min.X+x+1, bounds.Min.Y+height)
		segments[x] = sImg.SubImage(rect)
	}

	return segments, nil
}

// ImageToPNG takes an image and encodes it to PNG format, returning the byte data.
// It returns an error if the encoding operation fails.
func ImageToPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
