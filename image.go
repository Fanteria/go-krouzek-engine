package gke

import (
	"bytes"
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

func loadImage(path string) (*ebiten.Image, error) {
	image_buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(bytes.NewReader(image_buffer))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(image), nil
}

func loadImageToBlock(path string) (*Blok, error) {
	img, err := loadImage(path)
	if err != nil {
		return nil, err
	}
	return &Blok{image: img, coords: Coords{0, 0}, scale: Scale{1, 1}}, err
}

