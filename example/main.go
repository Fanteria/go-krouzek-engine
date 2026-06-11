package main

import (
	"embed"

	gke "github.com/Fanteria/go-krouzek-engine"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed all:obrazky
var assets embed.FS

func main() {
	gke.NastavSlozkuSObrazky(&assets)
	gke.NastavPozadi("./obrazky/pozadi.png")
	gke.NastavRezimPozadi(gke.RezimPozadiVyplnit)

	// Přidání statického bloku kamene
	blok := gke.PridejBlok("./obrazky/rock.png")
	gke.NastavZvetseni(blok, 0.55)
	gke.NastavPozici(blok, 20, 320)
	gke.NastavBlokovani(blok, true)

	// Přidání bloků podlahy
	v := gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32}
	gke.PridejBlokyPlosinuZBlokuSVyrezem("./obrazky/bloky.png", v, 15, 0, 400)
	gke.PridejBlokyPlosinuZBlokuSVyrezem("./obrazky/bloky.png", v, 5, 15*32+125, 350)
	gke.PridejBlokyPlosinuZBlokuSVyrezem("./obrazky/bloky.png", v, 10, (15+5)*32+125+150, 450)

	// Přidání animované postavy s různými druhy animace.
	hratelna_postava := gke.PrijdejHratelnouPostavu(
		"./obrazky/knight.png",
		0.1,
		map[ebiten.Key]gke.Akce{
			ebiten.KeyArrowRight: gke.AkceJdeVPravo,
			ebiten.KeyArrowLeft:  gke.AkceJdeVLevo,
			ebiten.KeySpace:      gke.AkceSkace,
		},
	)
	gke.NastavPozici(&hratelna_postava.Blok, 150.0, 333.0)
	gke.NastavZvetseni(&hratelna_postava.Blok, 3)
	gke.NastavRychlostPohybu(hratelna_postava, 3)
	gke.NastavSiluSkoku(hratelna_postava, 7)
	mriezka_postava := gke.NactiMriezku("./obrazky/knight.png")
	gke.NastavAnimaci(hratelna_postava, gke.AkceStoji, false, mriezka_postava.Snimky(0, 1, 2, 3))
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVPravo, false, mriezka_postava.Snimky(8, 9, 10, 11, 12, 13))
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVLevo, true, mriezka_postava.Snimky(8, 9, 10, 11, 12, 13))

	// Přidání animovaného bloku pomocí mřížky snímků.
	// NactiMriezku načte nastavení ze souboru tree_animated.png.mriezka.
	// Pokud soubor neexistuje, otevře se okno pro nastavení mřížky.
	ohen := gke.PridejAnimovanyBlok(
		"./obrazky/CampFireFinished.png",
		0.08,
		gke.NactiMriezku("./obrazky/CampFireFinished.png").VsechnySnimky()...,
	)
	gke.NastavPozici(ohen, 200, 360)
	gke.NastavZvetseni(ohen, 0.75)

	gke.NastavKameru(hratelna_postava)
	gke.NastavOkrajeKamery(200, 200, 100, 20)

	gke.SpustHru()
}
