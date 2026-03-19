package main

import (
	gke "github.com/Fanteria/go-krouzek-engine"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	gke.NastavPozadi("./obrazky/pozadi.png")

	blok := gke.PridejBlok("./obrazky/rock.png")
	gke.NastavZvetseni(blok, 0.5)
	gke.NastavPozici(blok, 50.0, 50.0)
	gke.NastavBlokovani(blok, true)

	blok_s_vyrezem := gke.PridejBlokSVyrezem("./obrazky/knight.png", gke.Vyrez{X1: 34, Y1: 5, X2: 49, Y2: 26})
	gke.NastavPozici(blok_s_vyrezem, 150.0, 150.0)

	animovany_blok := gke.PrijdejAnimovanyBlok(
		"./obrazky/knight.png",
		0.1,
		gke.Vyrez{X1: 34, Y1: 5, X2: 49, Y2: 24},
		gke.Vyrez{X1: 50, Y1: 5, X2: 65, Y2: 24},
		gke.Vyrez{X1: 66, Y1: 5, X2: 81, Y2: 24},
		gke.Vyrez{X1: 82, Y1: 5, X2: 97, Y2: 24},
	)
	gke.NastavPozici(animovany_blok, 250.0, 150.0)
	gke.NastavZvetseni(animovany_blok, 2)
	gke.NastavBlokovani(animovany_blok, true)

	hratelna_postava := gke.PrijdejHratelnouPostavu(
		"./obrazky/knight.png",
		0.1,
		map[ebiten.Key]gke.Akce{
			ebiten.KeyArrowRight: gke.AkceJdeVPravo,
			ebiten.KeyArrowLeft:  gke.AkceJdeVLevo,
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
