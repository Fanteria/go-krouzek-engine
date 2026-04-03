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
	AkceSkace
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

// TODO default should be zero in next year
var gravity float64 = 0.3

type Postava struct {
	Blok
	actionSubImages [AkcePocet]PostavaAnimation
	actualActions   []Akce
	animationSpeed  int
	speed           float64
	velocityY       float64
	jumpPower       float64
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
	var dy float64
	for _, action := range b.actualActions {
		switch action {
		case AkceJdeVPravo:
			dx += b.speed
		case AkceJdeVLevo:
			dx -= b.speed
		case AkceJdeNahoru:
			dy -= b.speed
		case AkceJdeDolu:
			dy += b.speed
		case AkceSkace:
			if b.velocityY == 0 {
				b.velocityY -= b.jumpPower
			}
		}
	}

	b.velocityY += gravity

	savedX := b.coords.x
	b.coords.x += dx
	if collidesWithSolid(b, blocks) {
		b.coords.x = savedX
	}

	savedY := b.coords.y
	if gravity == 0 {
		b.coords.y += dy
	}
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

type Enemy struct {
	Postava
	movingStrategy func(enemy *Enemy) []Akce
}

func (e *Enemy) move(blocks []drawable) {
	e.Postava.actualActions = e.movingStrategy(e)
	if len(e.Postava.actualActions) == 0 {
		e.Postava.actualActions = []Akce{AkceStoji}
	}
	e.Postava.move(blocks)
}
