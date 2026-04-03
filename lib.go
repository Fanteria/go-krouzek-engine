package gke

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO What to do with package gke problem in LSP?

// NastavUrovenLogovani nastaví, jak moc podrobné zprávy se budou vypisovat do terminálu.
// Můžeš použít: LogDebug (nejvíce zpráv), LogInfo, LogWarning, LogError (pouze chyby).
func NastavUrovenLogovani(uroven LogLevel) {
	setLogLevel(uroven)
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

func PrijdejAnimovanyBlok(cesta_k_obrazku string, rychlost_animace float64, vyrezy ...Vyrez) *Blok {
	// TODO to be removed for now leave this function here to not break API at now
	return PridejAnimovanyBlok(cesta_k_obrazku, rychlost_animace, vyrezy...)
}

// PridejAnimovanyBlok přidá do hry animovaný blok – obrázek, který se pohybuje jako animace.
// Zadej cestu k obrázku (spritesheetu), rychlost animace (např. 0.1 = pomalá, 1.0 = rychlá)
// a libovolný počet výřezů, které tvoří jednotlivé snímky animace.
// Vrátí ukazatel na blok.
func PridejAnimovanyBlok(cesta_k_obrazku string, rychlost_animace float64, vyrezy ...Vyrez) *Blok {
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek bloku", "chyba", err)
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

func PridejNepritele(cesta_k_obrazku string, strategie_pohybu func(*Enemy) []Akce) *Postava {
	// TODO rychlost_animace should be in separate function
	sub_block, err := loadImageToBlock(cesta_k_obrazku)
	if err != nil {
		log.Error("Nepodařilo se načíst obrázek nepřítele", "chyba", err)
		os.Exit(1)
	}
	block := &Enemy{
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

// NastavPozici přesune blok na zadané souřadnice na obrazovce.
// Souřadnice x udává vzdálenost od levého okraje, y od horního okraje (v pixelech).
func NastavPozici(blok *Blok, x float64, y float64) {
	blok.coords.x = x
	blok.coords.y = y
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

// SpustHru spustí hru! Tuhle funkci zavolej jako poslední, až budeš mít vše připraveno.
// Po jejím zavolání se otevře okno hry a hra začne běžet.
func SpustHru() {
	ebiten.RunGame(&game_instance)
}
