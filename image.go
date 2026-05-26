package gke

import (
	"bytes"
	"embed"
	"image"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

// Loaded images are in map so loadImage can be used multiple time without large memory overhead.
var images = make(map[string]*ebiten.Image)

// TODO doc comment
var assets *embed.FS = nil

// TODO doc comment
func readFile(path string) ([]byte, error) {
	if assets == nil {
		return os.ReadFile(path)
	} else {
		return assets.ReadFile(strings.TrimPrefix(path, "./"))
	}
}

// loadImage loads an image from the given file path and returns it as an ebiten image.
// Supports PNG and JPEG formats.
func loadImage(path string) (*ebiten.Image, error) {
	img := images[path]
	if img != nil {
		return img, nil
	}
	image_buffer, err := readFile(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(bytes.NewReader(image_buffer))
	if err != nil {
		return nil, err
	}
	img = ebiten.NewImageFromImage(image)
	images[path] = img
	return img, nil
}

// loadImageToBlock loads an image from the given file path and wraps it in a Blok
// with default position (0, 0) and scale (1, 1).
func loadImageToBlock(path string) (*Blok, error) {
	img, err := loadImage(path)
	if err != nil {
		return nil, err
	}
	return &Blok{image: img, coords: Coords{0, 0}, scale: Scale{1, 1}}, err
}
