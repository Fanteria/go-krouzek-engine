package gke

import "github.com/hajimehoshi/ebiten/v2"

// RezimPozadi určuje, jak se obrázek pozadí přizpůsobí velikosti okna.
type RezimPozadi int

const (
	// Roztahnout roztáhne obrázek přesně na velikost okna (může zkreslit poměr stran).
	RezimPozadiRoztahnout RezimPozadi = iota
	// Vyplnit zvětší obrázek tak, aby zaplnil celé okno – část obrázku může být oříznutá.
	RezimPozadiVyplnit
	// Prizpusobit zmenší nebo zvětší obrázek tak, aby byl celý vidět – mohou vzniknout okraje.
	RezimPozadiPrizpusobit
	// Puvodni zobrazí obrázek v jeho původní velikosti bez jakéhokoli škálování.
	RezimPozadiPuvodni
)

type background struct {
	image *ebiten.Image
	mode  RezimPozadi
}

// draw renders the background onto the screen using the configured RezimPozadi.
func (bg *background) draw(screen *ebiten.Image) {
	sw := float64(screen.Bounds().Dx())
	sh := float64(screen.Bounds().Dy())
	iw := float64(bg.image.Bounds().Dx())
	ih := float64(bg.image.Bounds().Dy())

	opts := &ebiten.DrawImageOptions{}

	switch bg.mode {
	case RezimPozadiRoztahnout:
		opts.GeoM.Scale(sw/iw, sh/ih)
	case RezimPozadiVyplnit:
		scale := max(sw/iw, sh/ih)
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate((sw-iw*scale)/2, (sh-ih*scale)/2)
	case RezimPozadiPrizpusobit:
		scale := min(sw/iw, sh/ih)
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate((sw-iw*scale)/2, (sh-ih*scale)/2)
	case RezimPozadiPuvodni:
		// no scaling
	}

	screen.DrawImage(bg.image, opts)
}
