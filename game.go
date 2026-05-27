package gke

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type game struct {
	outsideWidth   int
	outsideHeight  int
	animationIndex int
	background     *background
	blocks         []drawable
	movables       []movable
	camera         camera
	gridSpacing    float64
}

// Global instance of the game
var game_instance game = game{
	outsideWidth:   0,
	outsideHeight:  0,
	animationIndex: 0,
	background:     nil,
	blocks:         []drawable{},
	gridSpacing:    0,
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
		picture := b.image.SubImage(sub_image).(*ebiten.Image) 
		screen.DrawImage(picture, options)
	}

	// Draw coordinates grid
	if g.gridSpacing > 0 {
		startX := math.Floor(g.camera.offsetX/g.gridSpacing) * g.gridSpacing
		startY := math.Floor(g.camera.offsetY/g.gridSpacing) * g.gridSpacing

		for x := startX; x < g.camera.offsetX+float64(screen.Bounds().Dx()); x += g.gridSpacing {
			sx := x - g.camera.offsetX
			label := fmt.Sprintf("%.0f", x)

			text.Draw(screen, label, basicfont.Face7x13, int(sx)+2, 12, color.White)
			vector.StrokeLine(screen, float32(sx), 0, float32(sx), float32(screen.Bounds().Dy()), 1, color.Gray{Y: 100}, false)
		}

		for y := startY; y < g.camera.offsetY+float64(screen.Bounds().Dy()); y += g.gridSpacing {
			sy := y - g.camera.offsetY
			label := fmt.Sprintf("%.0f", y)

			text.Draw(screen, label, basicfont.Face7x13, 2, int(sy)+12, color.White)
			vector.StrokeLine(screen, 0, float32(sy), float32(screen.Bounds().Dx()), float32(sy), 1, color.Gray{Y: 100}, false)
		}
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	log.Debug("Nastavení velikosti okna hry", "šířka", outsideWidth, "výška", outsideHeight)
	return outsideWidth, outsideHeight
}
