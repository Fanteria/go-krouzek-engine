package main

import (
	"embed"
	"fmt"
	"os"

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

	// Přidání animovaného ohniště na konec poslední plošiny. To je cíl hráče.
	ohen := gke.PridejAnimovanyBlok(
		"./obrazky/CampFireFinished.png",
		0.08,
		gke.NactiMriezku("./obrazky/CampFireFinished.png").VsechnySnimky()...,
	)
	gke.NastavPozici(ohen, 1150, 420)
	gke.NastavZvetseni(ohen, 0.75)

	gke.NastavKameru(hratelna_postava)
	gke.NastavOkrajeKamery(200, 200, 100, 20)

	// Proměnné sledující stav hry – musí být deklarovány před uzávěrami, které je zachytávají.
	zivoty := 3
	kontaktCooldown := 0
	hadSmer := false // false = jde doleva, true = jde doprava
	fmt.Println(hadSmer)

	// Přidání hada (nepřítele), který hlídkuje na třetí plošině.
	// Had se pohybuje sem a tam mezi x=940 a x=1070.
	had := gke.PridejNepritele("./obrazky/characters.png", func(had *gke.Postava) []gke.Akce {
		if gke.ZjistitPoziciX(&had.Blok) > 1070 {
			hadSmer = false
		}
		if gke.ZjistitPoziciX(&had.Blok) < 940 {
			hadSmer = true
		}
		if hadSmer {
			return []gke.Akce{gke.AkceJdeVPravo}
		}
		return []gke.Akce{gke.AkceJdeVLevo}
	})
	gke.NastavPozici(&had.Blok, 970, 370)
	gke.NastavZvetseni(&had.Blok, 2)
	gke.NastavRychlostPohybu(had, 2)
	gke.NastavAnimaci(had, gke.AkceJdeVLevo, true, gke.NactiMriezku("./obrazky/characters.png").Rada(69, 73))
	gke.NastavAnimaci(had, gke.AkceJdeVPravo, false, gke.NactiMriezku("./obrazky/characters.png").Rada(69, 73))

	// Vytvoření dvou závěrečných obrazovek: výhra a smrt.
	// Tlačítko "hrát znovu" obnoví i proměnné stavu hry.
	obrazovkaVyhry := gke.NastavKonecovouObrazovku("Gratulujeme!", "Dohral jsi hru!")
	gke.PridejTlacitko(obrazovkaVyhry, "Hrat znovu", func() {
		zivoty = 3
		kontaktCooldown = 0
		hadSmer = false
		gke.ResetujHru()
	})
	gke.PridejTlacitko(obrazovkaVyhry, "Konec", func() { os.Exit(0) })

	obrazovkaSmerti := gke.NastavKonecovouObrazovku("Zemrel jsi!", "Zkus to znovu.")
	gke.PridejTlacitko(obrazovkaSmerti, "Zkusit znovu", func() {
		zivoty = 3
		kontaktCooldown = 0
		hadSmer = false
		gke.ResetujHru()
	})
	gke.PridejTlacitko(obrazovkaSmerti, "Konec", func() { os.Exit(0) })

	// Herní logika spouštěná každý snímek.
	gke.NastavAktualizaci(func() {
		// Kontakt s hadem: ubrat život s krátkým zpomalením, aby hráč nesepral přes víc životů najednou.
		if kontaktCooldown > 0 {
			kontaktCooldown--
		} else if gke.ZjistiKontaktSHratelnouPostavou(had) {
			zivoty--
			kontaktCooldown = 90
			if zivoty <= 0 {
				gke.ZobrazKonecovouObrazovku(obrazovkaSmerti)
			}
		}

		// Výhra: hráč dosáhl ohniště.
		if gke.ZjistitPoziciX(&hratelna_postava.Blok) >= gke.ZjistitPoziciX(ohen)-20 {
			gke.ZobrazKonecovouObrazovku(obrazovkaVyhry)
		}

		// Smrt: hráč vypadl ze světa.
		if gke.ZjistitPoziciY(&hratelna_postava.Blok) > 600 {
			gke.ZobrazKonecovouObrazovku(obrazovkaSmerti)
		}
	})

	gke.SpustHru()
}
