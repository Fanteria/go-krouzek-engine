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
	move()
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

type Postava struct {
	Blok
	actionSubImages [AkcePocet]PostavaAnimation
	actualActions   []Akce
	animationSpeed  int
	speed           float64
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

func (b *Postava) move() {
	for _, action := range b.actualActions {
		switch action {
		case AkceJdeVPravo:
			b.coords.x += b.speed
		case AkceJdeVLevo:
			b.coords.x -= b.speed
		case AkceJdeNahoru:
			b.coords.y -= b.speed
		case AkceJdeDolu:
			b.coords.y += b.speed
		}
	}
}

type HratelnaPostava struct {
	Postava
	moveActions map[ebiten.Key]Akce
}

func (p *HratelnaPostava) move() {
	var pressed_keys []ebiten.Key
	pressed_keys = inpututil.AppendPressedKeys(pressed_keys)
	p.Postava.actualActions = []Akce{}
	for _, pressed_key := range pressed_keys {
		action := p.moveActions[pressed_key]
		if action != AkceNic {
			p.Postava.actualActions = append(p.Postava.actualActions, action)
		}
	}
	p.Postava.move()
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
		movable.move()
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
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Scale(b.scale.width, b.scale.height)
		if block.isMirrored() {
			// TODO this will rotate it around axis on in place
			options.GeoM.Scale(-1, 1)
		}
		options.GeoM.Translate(b.coords.x, b.coords.y)
		sub_image := block.getSubImage(g.animationIndex)
		picture := ebiten.NewImageFromImage(b.image.SubImage(sub_image))
		screen.DrawImage(ebiten.NewImageFromImage(picture), options)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.outsideWidth = outsideWidth
	g.outsideHeight = outsideHeight
	return 320, 240
}
