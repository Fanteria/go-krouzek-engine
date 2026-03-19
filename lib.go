package gke

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func NastavPozadi(cesta_k_obrazku string) {
	img, err := loadImage(cesta_k_obrazku)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}
	game_instance.background = &Background{image: img}
}

func PridejBlok(cesta_k_obrazku string) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}
	block := &StatickyBlok{
		Blok:     *sub_block,
		subImage: sub_block.image.Bounds(),
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return &block.Blok
}

func PridejBlokSVyrezem(cesta_k_obrazku string, vyrez Vyrez) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}
	block := &StatickyBlok{
		Blok:     *sub_block,
		subImage: image.Rect(vyrez.X1, vyrez.Y1, vyrez.X2, vyrez.Y2),
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return &block.Blok
}

func PrijdejAnimovanyBlok(cesta_k_obrazku string, rychlost_animace float64, vyrezy ...Vyrez) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}
	var subImages []image.Rectangle
	for _, vyrez := range vyrezy {
		subImages = append(subImages, image.Rect(vyrez.X1, vyrez.Y1, vyrez.X2, vyrez.Y2))
	}
	animationSpeed := int(1 / rychlost_animace)
	if animationSpeed <= 0 {
		animationSpeed = 1
	}
	block := &AnimovanyBlok{
		Blok:           *sub_block,
		animationSpeed: animationSpeed,
		subImage:       subImages,
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return &block.Blok
}

func PrijdejHratelnouPostavu(cesta_k_obrazku string, rychlost_animace float64, akce_pohybu map[ebiten.Key]Akce) *Postava {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		fmt.Errorf("Error: %v", err)
		os.Exit(1)
	}
	animationSpeed := int(1 / rychlost_animace)
	if animationSpeed <= 0 {
		animationSpeed = 1
	}
	block := &HratelnaPostava{
		Postava: Postava{
			Blok:           *sub_block,
			actualActions:  []Akce{AkceStoji},
			animationSpeed: animationSpeed,
			speed:          1,
		},
		moveActions: akce_pohybu,
	}
	game_instance.blocks = append(game_instance.blocks, block)
	game_instance.movables = append(game_instance.movables, block)
	return &block.Postava
}

func NastavPevnost(blok *Blok) {
	blok.solid = true
}

func NastavZvetseni(blok *Blok, zvetseni float64) {
	blok.scale.width = zvetseni
	blok.scale.height = zvetseni
}

func NastavPozici(blok *Blok, x float64, y float64) {
	blok.coords.x = x
	blok.coords.y = y
}

func NastavBlokovani(blok *Blok, blokuje bool) {
	blok.solid = blokuje
}

func NastavRychlostPohybu(postava *Postava, rychlost_pohybu float64) {
	postava.speed = rychlost_pohybu
}

func NastavAnimaci(postava *Postava, akce Akce, zrcadlove_otocena bool, animace []Vyrez) {
	var rectangles []image.Rectangle
	for _, rect := range animace {
		rectangles = append(rectangles, image.Rect(rect.X1, rect.Y1, rect.X2, rect.Y2))
	}
	postava.actionSubImages[akce] = PostavaAnimation{
		mirror:     zrcadlove_otocena,
		rectangles: rectangles,
	}
}

func SpustHru() {
	ebiten.RunGame(&game_instance)
}
