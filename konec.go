package gke

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// KonecovyObrazovka je obrazovka zobrazená na konci hry – může obsahovat text a tlačítka.
type KonecovyObrazovka struct {
	textLines []string
	buttons   []tlacitko
}

type tlacitko struct {
	label  string
	action func()
}

const (
	konecLineHeight  = 20
	konecGap         = 36
	konecBtnHeight   = 28
	konecBtnPadding  = 10
	konecBtnMinWidth = 140
	konecBtnLabelPad = 20
	konecCharWidth   = 7
	konecBaseline    = 13
)

func (k *KonecovyObrazovka) layout(screenW, screenH int) (textStartY int, btnRects []image.Rectangle) {
	textH := len(k.textLines) * konecLineHeight
	buttonsH := len(k.buttons)*konecBtnHeight + max(0, len(k.buttons)-1)*konecBtnPadding
	totalH := textH + konecGap + buttonsH

	startY := (screenH - totalH) / 2
	if startY < 10 {
		startY = 10
	}
	textStartY = startY

	btnY := startY + textH + konecGap
	btnRects = make([]image.Rectangle, len(k.buttons))
	for i, btn := range k.buttons {
		btnW := len(btn.label)*konecCharWidth + 2*konecBtnLabelPad
		if btnW < konecBtnMinWidth {
			btnW = konecBtnMinWidth
		}
		x := (screenW - btnW) / 2
		btnRects[i] = image.Rect(x, btnY, x+btnW, btnY+konecBtnHeight)
		btnY += konecBtnHeight + konecBtnPadding
	}
	return
}

func (k *KonecovyObrazovka) draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	vector.DrawFilledRect(screen, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 180}, false)

	textStartY, btnRects := k.layout(w, h)
	mx, my := ebiten.CursorPosition()

	for i, line := range k.textLines {
		tw := len(line) * konecCharWidth
		tx := (w - tw) / 2
		ty := textStartY + i*konecLineHeight + konecBaseline
		text.Draw(screen, line, basicfont.Face7x13, tx, ty, color.White)
	}

	for i, btn := range k.buttons {
		r := btnRects[i]
		hovered := mx >= r.Min.X && mx < r.Max.X && my >= r.Min.Y && my < r.Max.Y
		bgColor := color.RGBA{60, 60, 80, 220}
		if hovered {
			bgColor = color.RGBA{100, 100, 140, 255}
		}
		vector.DrawFilledRect(screen,
			float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()),
			bgColor, false)
		vector.StrokeRect(screen,
			float32(r.Min.X), float32(r.Min.Y),
			float32(r.Dx()), float32(r.Dy()),
			1, color.RGBA{180, 180, 220, 255}, false)
		tw := len(btn.label) * konecCharWidth
		tx := r.Min.X + (r.Dx()-tw)/2
		ty := r.Min.Y + (r.Dy()+konecBaseline)/2 - 3
		text.Draw(screen, btn.label, basicfont.Face7x13, tx, ty, color.White)
	}
}

func (k *KonecovyObrazovka) checkClick(x, y, screenW, screenH int) func() {
	_, btnRects := k.layout(screenW, screenH)
	for i, r := range btnRects {
		if x >= r.Min.X && x < r.Max.X && y >= r.Min.Y && y < r.Max.Y {
			return k.buttons[i].action
		}
	}
	return nil
}
