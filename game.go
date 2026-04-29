package gke

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type game struct {
	outsideWidth   int
	outsideHeight  int
	animationIndex int
	background     *background
	blocks         []drawable
	movables       []movable
	camera         camera
}

// Global instance of the game
var game_instance game = game{
	outsideWidth:   0,
	outsideHeight:  0,
	animationIndex: 0,
	background:     nil,
	blocks:         []drawable{},
}

func (g *game) Update() error {
	g.animationIndex += 1
	for _, movable := range g.movables {
		movable.move(g.blocks)
	}
	g.camera.actualize(g.outsideWidth, g.outsideHeight)
	log.Debug("Camera actualized", "camera", g.camera)
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	var spacing float32 = 100.0
	w, h := screen.Size()
	for x := float32(0); x < float32(w); x += spacing {
		vector.StrokeLine(screen, x, 0, x, float32(h), 1, color.Gray{Y: 100}, false)
	}

	for y := float32(0); y < float32(h); y += spacing {
		vector.StrokeLine(screen, 0, y, float32(w), y, 1, color.Gray{Y: 100}, false)
	}

	// rendering background
	g.background.draw(screen, g.camera.offsetX, g.camera.offsetY)
	// rendering blocks
	for _, block := range g.blocks {
		b := block.getBlock()
		sub_image := block.getSubImage(g.animationIndex)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(b.scale.width, b.scale.height)
		if block.isMirrored() {
			w := float64(sub_image.Dx()) * b.scale.width
			options.GeoM.Scale(-1, 1)
			options.GeoM.Translate(w, 0)
		}
		options.GeoM.Translate(b.coords.x-g.camera.offsetX, b.coords.y-g.camera.offsetY)
		picture := ebiten.NewImageFromImage(b.image.SubImage(sub_image))
		screen.DrawImage(picture, options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	log.Debug("Nastavení velikosti okna hry", "šířka", outsideWidth, "výška", outsideHeight)
	return outsideWidth, outsideHeight
}
