package main

import (
	"embed"
	"fmt"

	gke "github.com/Fanteria/go-krouzek-engine"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed all:obrazky
var assets embed.FS

func main() {
	gke.NastavSlozkuSObrazky(&assets)
	gke.NastavUrovenLogovani(gke.LogWarning)
	// TODO fix name to NastavSouradnicovouMrizku
	gke.NastavSouradnicivouMrizku(50)
	gke.NastavGravitaci(0.3)
	gke.NastavPozadi("./obrazky/pozadi.png")
	gke.NastavRezimPozadi(gke.RezimPozadiVyplnit)

	// Přidání statického bloku kamene
	blok := gke.PridejBlok("./obrazky/rock.png")
	gke.NastavZvetseni(blok, 0.55)
	gke.NastavPozici(blok, 20, 320)
	gke.NastavBlokovani(blok, true)

	mine_stone := gke.PridejBlokSVyrezem(
		"./obrazky/mine-resources.png",
		gke.NactiMriezku("./obrazky/mine-resources.png").Snimky(13)[0])
	gke.NastavZvetseni(mine_stone, 3)
	gke.NastavPozici(mine_stone, 1100, 420)
	gke.NastavBlokovani(mine_stone, true)

	// Přidání bloků podlahy
	gke.PridejBlokyPlosinuZBlokuSVyrezem(
		"./obrazky/bloky.png",
		gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32},
		15,
		0, 400,
	)
	gke.PridejBlokyPlosinuZBlokuSVyrezem(
		"./obrazky/bloky.png",
		gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32},
		5,
		15*32+125, 350,
	)
	gke.PridejBlokyPlosinuZBlokuSVyrezem(
		"./obrazky/bloky.png",
		gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32},
		10,
		(15+5)*32+125+150, 450,
	)

	// Přidání animovaného bloku pomocí mřížky snímků.
	// NactiMriezku načte nastavení ze souboru tree_animated.png.mriezka.
	// Pokud soubor neexistuje, otevře se okno pro nastavení mřížky.
	// Pro tento obrázek nastav: Sirka snimku=64, Vyska snimku=64, Offset X=0, Offset Y=0.
	mriezka_strom := gke.NactiMriezku("./obrazky/tree_animated.png")
	animovany_blok := gke.PridejAnimovanyBlok(
		"./obrazky/tree_animated.png",
		0.08,
		mriezka_strom.VsechnySnimky()...,
	)
	gke.NastavPozici(animovany_blok, 250.0, 208.0)
	gke.NastavZvetseni(animovany_blok, 3)

	// Přidání animované postavy s různými druhy animace.
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
	gke.NastavPozici(&hratelna_postava.Blok, 150.0, 333.0)
	gke.NastavZvetseni(&hratelna_postava.Blok, 3)
	gke.NastavRychlostPohybu(hratelna_postava, 3)
	gke.NastavSiluSkoku(hratelna_postava, 7)
	mriezka_postava := gke.NactiMriezku("./obrazky/knight.png")
	gke.NastavAnimaci(hratelna_postava, gke.AkceStoji, false, mriezka_postava.Snimky(0, 1, 2, 3))
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVPravo, false, mriezka_postava.Snimky(8, 9, 10, 11, 12, 13))
	gke.NastavAnimaci(hratelna_postava, gke.AkceJdeVLevo, true, mriezka_postava.Snimky(8, 9, 10, 11, 12, 13))
	gke.NastavAnimaci(hratelna_postava, gke.AkceSkace, false,
		[]gke.Vyrez{
			{X1: 34, Y1: 47, X2: 49, Y2: 68},
			{X1: 50, Y1: 47, X2: 65, Y2: 68},
			{X1: 66, Y1: 47, X2: 81, Y2: 68},
			{X1: 82, Y1: 47, X2: 97, Y2: 68},
		},
	)
	gke.NastavAnimaci(hratelna_postava, gke.AkceSkaceVPravo, false,
		[]gke.Vyrez{
			{X1: 34, Y1: 47, X2: 49, Y2: 68},
			{X1: 50, Y1: 47, X2: 65, Y2: 68},
			{X1: 66, Y1: 47, X2: 81, Y2: 68},
			{X1: 82, Y1: 47, X2: 97, Y2: 68},
		},
	)
	gke.NastavAnimaci(hratelna_postava, gke.AkceSkaceVLevo, true,
		[]gke.Vyrez{
			{X1: 34, Y1: 47, X2: 49, Y2: 68},
			{X1: 50, Y1: 47, X2: 65, Y2: 68},
			{X1: 66, Y1: 47, X2: 81, Y2: 68},
			{X1: 82, Y1: 47, X2: 97, Y2: 68},
		},
	)

	var direction gke.Akce = gke.AkceJdeVPravo
	npc := gke.PridejNepritele(
		"./obrazky/characters.png",
		func(enemy *gke.Postava) []gke.Akce {
			x := gke.ZjistitPoziciX(&enemy.Blok)
			if x >= 450 {
				direction = gke.AkceJdeVLevo
			} else if x <= 250 {
				direction = gke.AkceJdeVPravo
			}
			return []gke.Akce{direction}
		},
	)
	gke.NastavPozici(&npc.Blok, 350.0, 305.0)
	gke.NastavBlokovani(&npc.Blok, false)
	gke.NastavZvetseni(&npc.Blok, 2.2)
	gke.NastavAnimaci(npc, gke.AkceJdeVPravo, false, gke.NactiMriezku("./obrazky/characters.png").Rada(50, 55))
	gke.NastavAnimaci(npc, gke.AkceJdeVLevo, true, gke.NactiMriezku("./obrazky/characters.png").Rada(50, 55))

	ohen := gke.PridejAnimovanyBlok(
		"./obrazky/CampFireFinished.png",
		0.08,
		gke.NactiMriezku("./obrazky/CampFireFinished.png").VsechnySnimky()...,
	)
	gke.NastavPozici(ohen, 200, 360)
	gke.NastavZvetseni(ohen, 0.75)

	zdravi := 100
	var snake_direction gke.Akce = gke.AkceStoji
	snake := gke.PridejNepritele(
		"./obrazky/characters.png",
		func(enemy *gke.Postava) []gke.Akce {
			if gke.ZjistiKontaktSHratelnouPostavou(enemy) {
				zdravi -= 1
				fmt.Println("Zdraví hráče:", zdravi)
			}
			return []gke.Akce{snake_direction}
		},
	)
	gke.NastavAnimaci(snake, gke.AkceStoji, true, gke.NactiMriezku("./obrazky/characters.png", "snake").Snimky(0, 2, 4, 6))
	gke.NastavPozici(&snake.Blok, 650, 300)
	gke.NastavZvetseni(&snake.Blok, 2.0)

	gke.NastavKameru(hratelna_postava)
	gke.NastavOkrajeKamery(200, 200, 100, 20)

	gke.SpustHru()
}
