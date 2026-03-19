package gke

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Akce int

const (
	AkceNic = iota
	AkceStoji
	AkceJdeVPravo
	AkceJdeVLevo
	AkceJdeNahoru
	AkceJdeDolu
	AkcePocet
)

type drawable interface {
	getSubImage(index int) image.Rectangle
	getBlock() Blok
	isMirrored() bool
}

type movable interface {
	// TODO there should be interface block or something like that
	move(blocks []drawable)
}

type Background struct {
	image *ebiten.Image
}

type Coords struct {
	x float64
	y float64
}

type Scale struct {
	width  float64
	height float64
}

type Vyrez struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

type Blok struct {
	image  *ebiten.Image
	coords Coords
	scale  Scale
	solid  bool
}

func (b *Blok) getBlock() Blok {
	return *b
}

func (b *Blok) getSubImage(index int) image.Rectangle {
	return b.image.Bounds()
}

func (b *Blok) isMirrored() bool {
	return false
}

type StatickyBlok struct {
	Blok
	subImage image.Rectangle
}

func (b *StatickyBlok) getSubImage(index int) image.Rectangle {
	return b.subImage
}

type AnimovanyBlok struct {
	Blok
	subImage       []image.Rectangle
	animationSpeed int
}

func (b *AnimovanyBlok) getSubImage(index int) image.Rectangle {
	return b.subImage[(index/b.animationSpeed)%len(b.subImage)]
}

type PostavaAnimation struct {
	mirror     bool
	rectangles []image.Rectangle
}

const gravity = 0.3

type Postava struct {
	Blok
	actionSubImages [AkcePocet]PostavaAnimation
	actualActions   []Akce
	animationSpeed  int
	speed           float64
	velocityY       float64
}

func (b *Postava) getSubImageAnimation() *PostavaAnimation {
	var actualAction Akce
	if len(b.actualActions) == 0 {
		actualAction = AkceStoji
	} else {
		actualAction = b.actualActions[0]
	}
	anim := &b.actionSubImages[actualAction]
	if len(anim.rectangles) == 0 {
		anim = &b.actionSubImages[AkceStoji]
	}
	return anim
}

func (b *Postava) getSubImage(index int) image.Rectangle {
	subimage := b.getSubImageAnimation()
	return subimage.rectangles[(index/b.animationSpeed)%len(subimage.rectangles)]
}

func (b *Postava) isMirrored() bool {
	subimage := b.getSubImageAnimation()
	return subimage.mirror
}

func (b *Postava) move(blocks []drawable) {
	// TODO use for better results if there is no only rectangles "github.com/solarlune/resolv"
	worldBounds := func(d drawable, animIndex int) (minX, minY, maxX, maxY float64) {
		b := d.getBlock()
		sub := d.getSubImage(animIndex)
		w := float64(sub.Dx()) * b.scale.width
		h := float64(sub.Dy()) * b.scale.height
		return b.coords.x, b.coords.y, b.coords.x + w, b.coords.y + h
	}

	collidesWithSolid := func(b *Postava, blocks []drawable) bool {
		// Always use the standing animation frame for a stable collision box.
		// Running/other animations may have different frame sizes which would cause
		// the box to grow into the ground and block horizontal movement.
		standingAnim := &b.actionSubImages[AkceStoji]
		if len(standingAnim.rectangles) == 0 {
			return false
		}
		sub := standingAnim.rectangles[0]
		w := float64(sub.Dx()) * b.scale.width
		h := float64(sub.Dy()) * b.scale.height
		pMinX, pMinY := b.coords.x, b.coords.y
		pMaxX, pMaxY := pMinX+w, pMinY+h

		for _, d := range blocks {
			if !d.getBlock().solid {
				continue
			}
			oMinX, oMinY, oMaxX, oMaxY := worldBounds(d, game_instance.animationIndex)
			if pMaxX > oMinX && pMinX < oMaxX && pMaxY > oMinY && pMinY < oMaxY {
				return true
			}
		}
		return false
	}

	var dx float64
	for _, action := range b.actualActions {
		switch action {
		case AkceJdeVPravo:
			dx += b.speed
		case AkceJdeVLevo:
			dx -= b.speed
		case AkceJdeDolu:
			if b.velocityY == 0 {
				b.velocityY -= 5.0
			}
			// case AkceJdeNahoru:
			// 	dy -= b.speed
			// case AkceJdeDolu:
			// 	dy += b.speed
		}
	}

	b.velocityY += gravity

	savedX := b.coords.x
	b.coords.x += dx
	if collidesWithSolid(b, blocks) {
		b.coords.x = savedX
	}

	savedY := b.coords.y
	b.coords.y += b.velocityY
	if collidesWithSolid(b, blocks) {
		b.coords.y = savedY
		b.velocityY = 0
	}
}

type HratelnaPostava struct {
	Postava
	moveActions map[ebiten.Key]Akce
}

func (p *HratelnaPostava) move(blocks []drawable) {
	var pressed_keys []ebiten.Key
	pressed_keys = inpututil.AppendPressedKeys(pressed_keys)
	p.Postava.actualActions = []Akce{}
	for _, pressed_key := range pressed_keys {
		action := p.moveActions[pressed_key]
		if action != AkceNic {
			p.Postava.actualActions = append(p.Postava.actualActions, action)
		}
	}
	p.Postava.move(blocks)
}

type game struct {
	outsideWidth   int
	outsideHeight  int
	animationIndex int
	background     *Background
	blocks         []drawable
	movables       []movable
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
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// rendering background
	{
		options := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.background.image, options)
	}
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
		options.GeoM.Translate(b.coords.x, b.coords.y)
		picture := ebiten.NewImageFromImage(b.image.SubImage(sub_image))
		screen.DrawImage(picture, options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return 320, 240
}
