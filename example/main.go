package main

import (
	gke "github.com/Fanteria/go-krouzek-engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	gke.NastavUrovenLogovani(gke.LogInfo)
	gke.NastavGravitaci(0.3)
	gke.NastavPozadi("./obrazky/pozadi.png")
	gke.NastavRezimPozadi(gke.RezimPozadiVyplnit)

	blok := gke.PridejBlok("./obrazky/rock.png")
	gke.NastavZvetseni(blok, 0.5)
	gke.NastavPozici(blok, 100.0, 150.0)
	gke.NastavBlokovani(blok, true)

	blok_cislo_n := 0.0
	for blok_cislo_n <= 20 {
		blok1 := gke.PridejBlokSVyrezem("./obrazky/bloky.png", gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32})
		gke.NastavPozici(blok1, blok_cislo_n*32, 400)
		gke.NastavBlokovani(blok1, true)
		blok_cislo_n += 1
	}

	blok_s_vyrezem := gke.PridejBlokSVyrezem("./obrazky/knight.png", gke.Vyrez{X1: 34, Y1: 5, X2: 49, Y2: 26})
	gke.NastavPozici(blok_s_vyrezem, 150.0, 150.0)

	var animace []gke.Vyrez
	for i := 0; i < 37; i += 1 {
		animace = append(animace, gke.Vyrez{X1: i * 64, Y1: 0, X2: (i + 1) * 64, Y2: 64})
	}
	animovany_blok := gke.PrijdejAnimovanyBlok(
		"./obrazky/tree_animated.png",
		0.08,
		animace...,
	)
	gke.NastavPozici(animovany_blok, 250.0, 208.0)
	gke.NastavZvetseni(animovany_blok, 3)
	gke.NastavBlokovani(animovany_blok, false)

	hratelna_postava := gke.PrijdejHratelnouPostavu(
		"./obrazky/knight.png",
		0.1,
		map[ebiten.Key]gke.Akce{
			ebiten.KeyArrowRight: gke.AkceJdeVPravo,
			ebiten.KeyArrowLeft:  gke.AkceJdeVLevo,
			ebiten.KeySpace:      gke.AkceSkace,
			ebiten.KeyArrowDown:  gke.AkceJdeDolu,
			ebiten.KeyArrowUp:    gke.AkceJdeNahoru,
		},
	)
	gke.NastavPozici(&hratelna_postava.Blok, 150.0, 50.0)
	gke.NastavZvetseni(&hratelna_postava.Blok, 1.5)
	gke.NastavRychlostPohybu(hratelna_postava, 2)
	gke.NastavAnimaci(hratelna_postava, gke.AkceStoji, false,
		[]gke.Vyrez{
			{X1: 34, Y1: 5, X2: 49, Y2: 24},
			{X1: 50, Y1: 5, X2: 65, Y2: 24},
			{X1: 66, Y1: 5, X2: 81, Y2: 24},
			{X1: 82, Y1: 5, X2: 97, Y2: 24},
		},
	)
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVPravo, false,
		[]gke.Vyrez{
			{X1: 34, Y1: 25, X2: 49, Y2: 46},
			{X1: 50, Y1: 25, X2: 65, Y2: 46},
			{X1: 66, Y1: 25, X2: 81, Y2: 46},
			{X1: 82, Y1: 25, X2: 97, Y2: 46},
			{X1: 98, Y1: 25, X2: 113, Y2: 46},
			{X1: 114, Y1: 25, X2: 129, Y2: 46},
		},
	)
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVLevo, true,
		[]gke.Vyrez{
			{X1: 34, Y1: 25, X2: 49, Y2: 46},
			{X1: 50, Y1: 25, X2: 65, Y2: 46},
			{X1: 66, Y1: 25, X2: 81, Y2: 46},
			{X1: 82, Y1: 25, X2: 97, Y2: 46},
			{X1: 98, Y1: 25, X2: 113, Y2: 46},
			{X1: 114, Y1: 25, X2: 129, Y2: 46},
		},
	)

	gke.SpustHru()
}
