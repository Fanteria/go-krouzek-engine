# Engine pro kroužek GoCode

Beware: using this library is not a good idea. It is designed to be used in a course for children so it does not match best practices, implement things the right way or even have consistent naming. Functions used in the course are in Czech, others are in English just for my own convenience. Moreover this library does not have a stable API, it is changed every year during summer holidays to reflect the needs of the course.

Knihovna pro tvorbu her v jazyce Go.

## Jak začít

Každý program s touto knihovnou vypadá přibližně takto:

```go
package main

import (
	gke "github.com/Fanteria/go-krouzek-engine"
)


func main() {
    gke.NastavPozadi("pozadi.png")

    zem := gke.PridejBlok("zem.png")
    gke.NastavPozici(zem, 0, 400)
    gke.NastavBlokovani(zem, true)

    postava := gke.PridejHratelnouPostavu("hrdina.png", 0.1, map[ebiten.Key]gke.Akce{
        ebiten.KeyArrowLeft:  gke.AkceJdeVlevo,
        ebiten.KeyArrowRight: gke.AkceJdePravo,
    })

    gke.SpustHru() // tohle volej vždy jako poslední!
}
```

## Přehled funkcí

### Pozadí

#### `NastavPozadi(cesta_k_obrazku)`
Nastaví obrázek, který se zobrazí jako pozadí celé hry – za vším ostatním.

```go
gke.NastavPozadi("obloha.png")
```

### Bloky

Blok je obrázek, který se zobrazí na obrazovce. Může to být kámen, plošina, strom – cokoliv.

#### `PridejBlok(cesta_k_obrazku) *Blok`
Přidá do hry jednoduchý blok z obrázku.

```go
kamen := gke.PridejBlok("kamen.png")
```

#### `PridejBlokSVyrezem(cesta_k_obrazku, vyrez) *Blok`
Přidá blok, ale zobrazí jen část obrázku. Hodí se pro tzv. **spritesheety** – velké obrázky
obsahující více věcí najednou. Výřez zadáš souřadnicemi levého horního rohu (X1, Y1)
a pravého dolního rohu (X2, Y2).

```go
// Z obrázku vezme jen část od pixelu (0,0) do (32,32)
strom := gke.PridejBlokSVyrezem("stromy.png", gke.Vyrez{X1: 0, Y1: 0, X2: 32, Y2: 32})
```

#### `PridejAnimovanyBlok(cesta_k_obrazku, rychlost_animace, vyrezy...) *Blok`
Přidá blok s animací – obrázek, který se střídavě mění jako animovaný film.
Zadáš rychlost (číslo od 0 do 1, větší = rychlejší) a jednotlivé snímky animace jako výřezy.

```go
ohne := gke.PridejAnimovanyBlok("ohne.png", 0.2,
    gke.Vyrez{X1: 0,  Y1: 0, X2: 32, Y2: 32},
    gke.Vyrez{X1: 32, Y1: 0, X2: 64, Y2: 32},
    gke.Vyrez{X1: 64, Y1: 0, X2: 96, Y2: 32},
)
```

### Nastavení bloků

#### `NastavPozici(blok, x, y)`
Přesune blok na zadané místo na obrazovce. `x` je vzdálenost od levého okraje,
`y` od horního okraje (v pixelech). Levý horní roh obrazovky je (0, 0).

```go
gke.NastavPozici(kamen, 100, 300)
```

#### `NastavZvetseni(blok, zvetseni)`
Změní velikost bloku. `1.0` = původní velikost, `2.0` = dvakrát větší, `0.5` = poloviční.

```go
gke.NastavZvetseni(kamen, 2.0)
```

#### `NastavBlokovani(blok, blokuje)`
Zapíná (`true`) nebo vypíná (`false`) pevnost bloku. Flexibilnější verze `NastavPevnost`.

```go
gke.NastavBlokovani(dvere, false) // postavy projdou skrz
gke.NastavBlokovani(dvere, true)  // postavy se zastaví
```

### Postavy

Postava je speciální blok, který může hráč ovládat klávesnicí.

#### `PridejHratelnouPostavu(cesta_k_obrazku, rychlost_animace, akce_pohybu) *Postava`
Přidá postavu ovládanou hráčem. Zadáš obrázek (spritesheet), rychlost animace
a mapu kláves – každé klávese přiřadíš akci, která se má stát po jejím stisku.

```go
postava := gke.PridejHratelnouPostavu("hrdina.png", 0.1, map[ebiten.Key]gke.Akce{
    ebiten.KeyArrowLeft:  gke.AkceJdeVlevo,
    ebiten.KeyArrowRight: gke.AkceJdePravo,
    ebiten.KeySpace:      gke.AkceSkace,
})
```

#### `NastavRychlostPohybu(postava, rychlost)`
Nastaví, jak rychle se postava pohybuje. Výchozí hodnota je `1.0`.

```go
gke.NastavRychlostPohybu(postava, 3.0) // třikrát rychlejší
```

#### `NastavAnimaci(postava, akce, zrcadlove_otocena, animace)`
Přiřadí postavě sadu snímků (animaci) pro konkrétní akci. Parametr `zrcadlove_otocena`
způsobí, že se obrázek překlopí zrcadlově – hodí se např. pro pohyb doleva,
kdy použiješ stejné snímky jako pro pohyb doprava, jen otočené.

```go
// Animace běhu doprava
gke.NastavAnimaci(postava, gke.AkceJdePravo, false, []gke.Vyrez{
    {X1: 0,  Y1: 0, X2: 32, Y2: 32},
    {X1: 32, Y1: 0, X2: 64, Y2: 32},
})

// Animace běhu doleva – stejné snímky, jen zrcadlově otočené
gke.NastavAnimaci(postava, gke.AkceJdeVlevo, true, []gke.Vyrez{
    {X1: 0,  Y1: 0, X2: 32, Y2: 32},
    {X1: 32, Y1: 0, X2: 64, Y2: 32},
})
```

### Spuštění hry

#### `SpustHru()`
Spustí hru. Tuhle funkci volej vždy jako **úplně poslední** – až budeš mít vše připraveno.

```go
gke.SpustHru()
```

## Souřadnice na obrazovce

Možná tě překvapí, že osa Y jde **dolů**, ne nahoru jako v matematice.
Levý horní roh obrazovky je bod (0, 0).

```
(0,0) ──────────────────► X
  │
  │
  │
  ▼
  Y
```

Takže pokud chceš dát blok dolů na obrazovce, použij **větší** hodnotu Y.

## Tipy

- Vždy volej `SpustHru()` jako poslední.
- Bloky přidávej v pořadí, v jakém chceš, aby se kreslily – první přidaný blok bude nejníže (překryjí ho ostatní).
- Pro pevné bloky (zem, zdi) nezapomeň zavolat `NastavPevnost`.
