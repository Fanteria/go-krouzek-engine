package gke

import (
	"embed"
	"image"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO What to do with package gke problem in LSP?

// NastavUrovenLogovani nastaví, jak moc podrobné zprávy se budou vypisovat do terminálu.
// Můžeš použít: LogDebug (nejvíce zpráv), LogInfo, LogWarning, LogError (pouze chyby).
func NastavUrovenLogovani(uroven LogLevel) {
	setLogLevel(uroven)
}

// NastavSouradnicivouMrizku zapne zobrazení souřadnicové mřížky s daným rozestupem mezi čarami (v pixelech).
func NastavSouradnicivouMrizku(rozestup float64) {
	game_instance.gridSpacing = rozestup
}

// Nastaví složku s obrázky, která je vložená přímo do programu pomocí go:embed.
// Zavolej tuto funkci na začátku programu, pokud chceš obrázky zabalit do spustitelného souboru.
func NastavSlozkuSObrazky(slozka *embed.FS) {
	assets = slozka
}

// NastavGravitaci nastaví sílu gravitace – jak rychle padají postavy dolů.
// Větší číslo = silnější gravitace, 0 = žádná gravitace.
func NastavGravitaci(gravitace float64) {
	gravity = gravitace
}

// NastavPozadi nastaví obrázek pozadí hry.
// Zadej cestu k obrázku, který chceš použít jako pozadí – například "pozadi.png".
// Pozadí se zobrazí za všemi ostatními věcmi ve hře.
func NastavPozadi(cesta_k_obrazku string) {
	img, err := loadImage(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst pozadí", "chyba", err)
		os.Exit(1)
	}
	game_instance.background = &background{image: img, mode: RezimPozadiPuvodni}
}

// NastavRezimPozadi nastaví, jak se obrázek pozadí přizpůsobí velikosti okna.
// Můžeš použít: Roztahnout, Vyplnit, Prizpusobit nebo Puvodni.
func NastavRezimPozadi(rezim RezimPozadi) {
	game_instance.background.mode = rezim
}

// PridejBlok přidá do hry nový blok (obrázek).
// Zadej cestu k obrázku bloku – například "kamen.png".
// Blok se zobrazí na obrazovce a můžeš s ním dále pracovat pomocí dalších funkcí.
// Vrátí ukazatel na blok, který pak můžeš předávat jiným funkcím.
func PridejBlok(cesta_k_obrazku string) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek bloku", "chyba", err)
		os.Exit(1)
	}
	block := &StatickyBlok{
		Blok:     *sub_block,
		subImage: sub_block.image.Bounds(),
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return &block.Blok
}

// PridejBlokSVyrezem přidá do hry blok, ale zobrazí jen část obrázku – výřez.
// Hodí se, když máš jeden velký obrázek s více věcmi (tzv. spritesheet) a chceš
// vybrat jen jednu z nich. Výřez zadáš jako souřadnice rohů obdélníku (X1, Y1) a (X2, Y2).
// Vrátí ukazatel na blok.
func PridejBlokSVyrezem(cesta_k_obrazku string, vyrez Vyrez) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek bloku", "chyba", err)
		os.Exit(1)
	}
	block := &StatickyBlok{
		Blok:     *sub_block,
		subImage: image.Rect(vyrez.X1, vyrez.Y1, vyrez.X2, vyrez.Y2),
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return &block.Blok
}

func PrijdejAnimovanyBlok(cesta_k_obrazku string, vyrezy ...Vyrez) *Blok {
	// TODO to be removed for now leave this function here to not break API at now
	return PridejAnimovanyBlok(cesta_k_obrazku, 0.1, vyrezy...)
}

// PridejAnimovanyBlok přidá do hry animovaný blok – obrázek, který se pohybuje jako animace.
// Zadej cestu k obrázku (spritesheetu), rychlost animace (např. 0.1 = pomalá, 1.0 = rychlá)
// a libovolný počet výřezů, které tvoří jednotlivé snímky animace.
// Vrátí ukazatel na blok.
func PridejAnimovanyBlok(cesta_k_obrazku string, rychlost_animace float64, vyrezy ...Vyrez) *Blok {
	block := PridejAnimovanyBlokAnim(cesta_k_obrazku, vyrezy...)
	NastavRychlostAnimaceProAnimovanyBlok(block, rychlost_animace)
	return &block.Blok
}

// PridejAnimovanyBlok přidá do hry animovaný blok – obrázek, který se pohybuje jako animace.
// Zadej cestu k obrázku (spritesheetu), rychlost animace (např. 0.1 = pomalá, 1.0 = rychlá)
// a libovolný počet výřezů, které tvoří jednotlivé snímky animace.
// Vrátí ukazatel na blok.
func PridejAnimovanyBlokAnim(cesta_k_obrazku string, vyrezy ...Vyrez) *AnimovanyBlok {
	// TODO this function should be mainly used in next school year
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek bloku", "chyba", err)
		os.Exit(1)
	}
	var subImages []image.Rectangle
	for _, vyrez := range vyrezy {
		subImages = append(subImages, image.Rect(vyrez.X1, vyrez.Y1, vyrez.X2, vyrez.Y2))
	}
	block := &AnimovanyBlok{
		Blok:           *sub_block,
		animationSpeed: 10,
		subImage:       subImages,
	}
	game_instance.blocks = append(game_instance.blocks, block)
	return block
}

// TODO doc comment
func PridejAnimovanyBlokVsechnySnimky(cesta_k_obrazku string) *AnimovanyBlok {
	return PridejAnimovanyBlokAnim(cesta_k_obrazku, NactiMriezku(cesta_k_obrazku).VsechnySnimky()...)
}

// TODO doc comment
func PridejAnimovanyBlokRada(cesta_k_obrazku string, zacatek, konec int) *AnimovanyBlok {
	return PridejAnimovanyBlokAnim(cesta_k_obrazku, NactiMriezku(cesta_k_obrazku).Rada(zacatek, konec)...)
}

// TODO doc comment
func PridejAnimovanyBlokSnimky(cesta_k_obrazku string, indexy ...int) *AnimovanyBlok {
	return PridejAnimovanyBlokAnim(cesta_k_obrazku, NactiMriezku(cesta_k_obrazku).Snimky(indexy...)...)
}

func NastavRychlostAnimaceProAnimovanyBlok(animovany_blok *AnimovanyBlok, rychlost_animace float64) {
	animationSpeed := int(1 / rychlost_animace)
	if animationSpeed <= 0 {
		animationSpeed = 1
	}
	animovany_blok.animationSpeed = animationSpeed
}

func PrijdejHratelnouPostavu(cesta_k_obrazku string, rychlost_animace float64, akce_pohybu map[ebiten.Key]Akce) *Postava {
	// TODO misspell delete for next school year
	return PridejHratelnouPostavu(cesta_k_obrazku, rychlost_animace, akce_pohybu)
}

// PridejHratelnouPostavu přidá do hry postavu, kterou může hráč ovládat klávesnicí.
// Zadej cestu k obrázku postavy, rychlost animace a mapu kláves na akce pohybu –
// například klávesa šipka doleva způsobí, že postava půjde doleva.
// Vrátí ukazatel na postavu, se kterou pak můžeš pracovat pomocí dalších funkcí.
func PridejHratelnouPostavu(cesta_k_obrazku string, rychlost_animace float64, akce_pohybu map[ebiten.Key]Akce) *Postava {
	// TODO rychlost_animace should be in separate function
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek hratelné postavy", "chyba", err)
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
			speed:          1.0,
			velocityY:      0.0,
			jumpPower:      5.0,
		},
		moveActions: akce_pohybu,
	}
	game_instance.blocks = append(game_instance.blocks, block)
	game_instance.movables = append(game_instance.movables, block)
	return &block.Postava
}

func PridejNepritele(cesta_k_obrazku string, strategie_pohybu func(*Postava) []Akce) *Postava {
	// TODO rychlost_animace should be in separate function
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek nepřítele", "chyba", err)
		os.Exit(1)
	}
	block := &enemy{
		Postava: Postava{
			Blok:           *sub_block,
			actualActions:  []Akce{AkceStoji},
			animationSpeed: 10,
			speed:          1.0,
			velocityY:      0.0,
			jumpPower:      5.0,
		},
		movingStrategy: strategie_pohybu,
	}
	game_instance.blocks = append(game_instance.blocks, block)
	game_instance.movables = append(game_instance.movables, block)

	return &block.Postava
}

// NastavZvetseni změní velikost bloku.
// Hodnota 1.0 znamená původní velikost, 2.0 znamená dvakrát větší, 0.5 znamená poloviční.
func NastavZvetseni(blok *Blok, zvetseni float64) {
	blok.scale.width = zvetseni
	blok.scale.height = zvetseni
}

// TODO func NastavVelikostObrazku

// NastavPozici přesune blok na zadané souřadnice na obrazovce.
// Souřadnice x udává vzdálenost od levého okraje, y od horního okraje (v pixelech).
func NastavPozici(blok *Blok, x float64, y float64) {
	blok.coords.x = x
	blok.coords.y = y
}

// ZjistitKontakt zjistí, zda se dvě postavy navzájem dotýkají (překrývají).
// Vrátí true, pokud se postavy překrývají, jinak false.
func ZjistitKontakt(a *Postava, b *Postava) bool {
	bounds := func(p *Postava) (minX, minY, maxX, maxY float64) {
		sub, found := firstAnimationFrame(p)
		var w, h float64
		if found {
			w = float64(sub.Dx()) * p.scale.width
			h = float64(sub.Dy()) * p.scale.height
		} else {
			w = float64(p.image.Bounds().Dx()) * p.scale.width
			h = float64(p.image.Bounds().Dy()) * p.scale.height
		}
		return p.coords.x, p.coords.y, p.coords.x + w, p.coords.y + h
	}
	aMinX, aMinY, aMaxX, aMaxY := bounds(a)
	bMinX, bMinY, bMaxX, bMaxY := bounds(b)
	return aMaxX > bMinX && aMinX < bMaxX && aMaxY > bMinY && aMinY < bMaxY
}

func ZjistiKontaktSHratelnouPostavou(postava *Postava) bool {
	return ZjistitKontakt(game_instance.camera.character, postava)
}

// ZjistitPoziciX vrátí aktuální souřadnici X (vodorovnou polohu) bloku v pixelech.
func ZjistitPoziciX(blok *Blok) float64 {
	return blok.coords.x
}

// ZjistitPoziciY vrátí aktuální souřadnici Y (svislou polohu) bloku v pixelech.
func ZjistitPoziciY(blok *Blok) float64 {
	return blok.coords.y
}

// NastavBlokovani zapíná nebo vypíná, zda blok zastavuje postavy.
// Pokud zadáš true, postavy se o blok zastaví. Pokud false, postavy projdou skrz.
func NastavBlokovani(blok *Blok, blokuje bool) {
	blok.solid = blokuje
}

// NastavRychlostPohybu nastaví, jak rychle se postava pohybuje.
// Výchozí rychlost je 1.0. Větší číslo = rychlejší pohyb, menší číslo = pomalejší pohyb.
func NastavRychlostPohybu(postava *Postava, rychlost_pohybu float64) {
	postava.speed = rychlost_pohybu
}

// NastavSiluSkoku nastaví, jak vysoko postava skáče.
// Větší číslo znamená vyšší skok.
func NastavSiluSkoku(postava *Postava, sila_skoku float64) {
	postava.jumpPower = sila_skoku
}

// NastavAnimaci přiřadí postavě animaci pro určitou akci (např. běh, stání, skok).
// Zadej postavu, akci, pro kterou animaci nastavuješ, zda má být obrázek zrcadlově otočený,
// a seznam výřezů ze spritesheetu, které tvoří snímky animace.
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

func NastavRychlostAnimace(postava *Postava, rychlost_animace float64) {
	animationSpeed := int(1 / rychlost_animace)
	if animationSpeed <= 0 {
		log.Warn("Rychlost animace nemůže být záportná", "rychlost animace", rychlost_animace)
		animationSpeed = 1
	}
}

// NastavKameru zapne sledování postavy kamerou – obrazovka se bude posouvat spolu s postavou.
// NastavKameru zapne sledování postavy kamerou – obrazovka se bude posouvat spolu s postavou.
// Kamera se nezasune doleva za souřadnici 0.
// Jak daleko od okrajů se kamera začne pohybovat nastavíš pomocí NastavOkrajeKamery.
func NastavKameru(postava *Postava) {
	game_instance.camera.character = postava
	game_instance.camera.active = true
}

// NastavOkrajeKamery nastaví, jak daleko od okrajů obrazovky musí postava být,
// aby se kamera začala posouvat. Například hodnoty 200 a 200 a 150 a 150 znamenají,
// že postava se může pohybovat v prostřední části obrazovky a kamera se pohne teprve
// když dojde blíž k některému okraji.
func NastavOkrajeKamery(vlevo, vpravo, nahoru, dolu float64) {
	game_instance.camera.marginLeft = vlevo
	game_instance.camera.marginRight = vpravo
	game_instance.camera.marginUp = nahoru
	game_instance.camera.marginDown = dolu
}

// ZapniKameru zapne posouvání obrazovky za postavou.
func ZapniKameru() {
	game_instance.camera.active = true
}

// VypniKameru vypne posouvání obrazovky – kamera zůstane na místě.
func VypniKameru() {
	game_instance.camera.active = false
}

func KameraOffsetX() float64 {
	return game_instance.camera.offsetX
}

func KameraOffsetY() float64 {
	return game_instance.camera.offsetY
}

// NactiMriezku načte mřížku snímků pro obrázek na dané cestě.
// Volitelný parametr nazev umožňuje mít více různých mřížek pro jeden obrázek –
// například gke.NactiMriezku("./postava.png", "stani") a gke.NactiMriezku("./postava.png", "beh").
// Pokud soubor s nastavením ještě neexistuje, otevře se okno pro konfiguraci mřížky.
// Po nastavení se konfigurace uloží – příště se načte automaticky.
func NactiMriezku(cesta_k_obrazku string, nazev ...string) *Mriezka {
	name := ""
	if len(nazev) > 0 {
		name = nazev[0]
	}
	if m, ok := nactiMriezku(cesta_k_obrazku, name); ok {
		return m
	}
	m, err := openGridConfigurator(cesta_k_obrazku, name)
	if err != nil {
		// TODO update log
		log.Error("Nepodařilo se načíst obrázek mřížky", "chyba", err)
		os.Exit(1)
	}
	return m
}

func UpravMriezku(cesta_k_obrazku string, nazev ...string) *Mriezka {
	name := ""
	if len(nazev) > 0 {
		name = nazev[0]
	}
	m, err := openGridConfigurator(cesta_k_obrazku, name)
	if err != nil {
		// TODO update log
		log.Error("Nepodařilo se načíst obrázek mřížky", "chyba", err)
		os.Exit(1)
	}
	return m
}

func PridejBlokyPlosinuZBloku(cesta_k_obrazku string, opakovani int, x, y float64) {
	for i := 0; i <= opakovani; i += 1 {
		blok := PridejBlok("./obrazky/bloky.png")
		blok.coords.x = x + float64(i*blok.image.Bounds().Dx())
		blok.coords.y = y
		blok.solid = true
	}
}

func PridejBlokyPlosinuZBlokuSVyrezem(cesta_k_obrazku string, vyrez Vyrez, opakovani int, x, y float64) {
	for i := 0; i <= opakovani; i += 1 {
		blok := PridejBlokSVyrezem(cesta_k_obrazku, vyrez)
		blok.coords.x = x + float64(i)*math.Abs(float64(vyrez.X1-vyrez.X2))
		blok.coords.y = y
		blok.solid = true
	}
}

// NastavKonecovouObrazovku vytvoří novou závěrečnou obrazovku s řádky textu.
// Můžeš zadat jeden nebo více řádků textu, například: "Gratulujeme!", "Dohrál jsi hru!".
// Vrátí ukazatel na obrazovku, ke které pak přidáš tlačítka pomocí PridejTlacitko.
func NastavKonecovouObrazovku(radky_textu ...string) *KonecovyObrazovka {
	return &KonecovyObrazovka{textLines: radky_textu}
}

// PridejTlacitko přidá tlačítko na závěrečnou obrazovku.
// Zadej obrazovku, popis tlačítka (text na tlačítku) a funkci, která se spustí po kliknutí.
func PridejTlacitko(obrazovka *KonecovyObrazovka, popis string, akce func()) {
	obrazovka.buttons = append(obrazovka.buttons, tlacitko{label: popis, action: akce})
}

// ZobrazKonecovouObrazovku zobrazí závěrečnou obrazovku přes hru a pozastaví herní smyčku.
func ZobrazKonecovouObrazovku(obrazovka *KonecovyObrazovka) {
	game_instance.endScreen = obrazovka
}

// SkryjKonecovouObrazovku skryje závěrečnou obrazovku a obnoví běh hry.
func SkryjKonecovouObrazovku() {
	game_instance.endScreen = nil
}

// ResetujHru obnoví hru do stavu před jejím spuštěním.
// Všechny postavy a bloky se vrátí na svá původní místa.
func ResetujHru() {
	game_instance.restoreSnapshot()
	game_instance.endScreen = nil
}

// NastavAktualizaci nastaví funkci, která se zavolá každý snímek hry.
// Použij ji pro vlastní herní logiku – například kontrolu, zda postava dosáhla cíle.
func NastavAktualizaci(callback func()) {
	game_instance.updateCallback = callback
}

// SpustHru spustí hru! Tuhle funkci zavolej jako poslední, až budeš mít vše připraveno.
// Po jejím zavolání se otevře okno hry a hra začne běžet.
func SpustHru() {
	game_instance.saveSnapshot()
	ebiten.RunGame(&game_instance)
}
