package gke

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"os"
	"strconv"
	"unicode"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Mriezka popisuje rovnoměrnou mřížku snímků ve spritesheetu.
type Mriezka struct {
	name        string
	imagePath   string
	imageWidth  int
	imageHeight int
	frameWidth  int
	frameHeight int
	offsetX     int
	offsetY     int
}

type mriezkaJSON struct {
	FrameWidth  int `json:"frameWidth"`
	FrameHeight int `json:"frameHeight"`
	OffsetX     int `json:"offsetX"`
	OffsetY     int `json:"offsetY"`
}

// Snimky vrátí výřezy pro zadané indexy snímků (číslované od 0).
func (m *Mriezka) Snimky(indexy ...int) []Vyrez {
	framesPerRow := 1
	if m.frameWidth > 0 {
		n := (m.imageWidth - m.offsetX) / m.frameWidth
		if n > 0 {
			framesPerRow = n
		}
	}
	result := make([]Vyrez, 0, len(indexy))
	for _, idx := range indexy {
		col := idx % framesPerRow
		row := idx / framesPerRow
		x1 := m.offsetX + col*m.frameWidth
		y1 := m.offsetY + row*m.frameHeight
		result = append(result, Vyrez{
			X1: x1,
			Y1: y1,
			X2: x1 + m.frameWidth,
			Y2: y1 + m.frameHeight,
		})
	}
	return result
}

// Rada vrátí výřezy pro snímky v rozsahu [zacatek, konec).
func (m *Mriezka) Rada(zacatek, konec int) []Vyrez {
	if konec <= zacatek {
		return nil
	}
	indexy := make([]int, konec-zacatek)
	for i := range indexy {
		indexy[i] = zacatek + i
	}
	return m.Snimky(indexy...)
}

// VsechnySnimky vrátí výřezy pro všechny snímky v mřížce.
func (m *Mriezka) VsechnySnimky() []Vyrez {
	if m.frameWidth <= 0 || m.frameHeight <= 0 {
		return nil
	}
	cols := (m.imageWidth - m.offsetX) / m.frameWidth
	rows := (m.imageHeight - m.offsetY) / m.frameHeight
	indexy := make([]int, cols*rows)
	for i := range indexy {
		indexy[i] = i
	}
	return m.Snimky(indexy...)
}

// --- config UI ---

const (
	konfFieldSirka   = 0
	konfFieldVyska   = 1
	konfFieldOffsetX = 2
	konfFieldOffsetY = 3
)

var konfFieldLabels = [4]string{
	"Sirka snimku",
	"Vyska snimku",
	"Offset X",
	"Offset Y",
}

type mriezkaKonfig struct {
	mriezka      *Mriezka
	rawImg       image.Image
	img          *ebiten.Image // created lazily inside RunGame session
	fields       [4]string
	active       int
	screenWidth  int
	screenHeight int
}

func (k *mriezkaKonfig) Update() error {
	var chars []rune
	chars = ebiten.AppendInputChars(chars)
	for _, ch := range chars {
		if unicode.IsDigit(ch) {
			k.fields[k.active] += string(ch)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(k.fields[k.active]) > 0 {
		r := []rune(k.fields[k.active])
		k.fields[k.active] = string(r[:len(r)-1])
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			k.active = (k.active + 3) % 4
		} else {
			k.active = (k.active + 1) % 4
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		fw, err1 := strconv.Atoi(k.fields[konfFieldSirka])
		fh, err2 := strconv.Atoi(k.fields[konfFieldVyska])
		if err1 != nil || err2 != nil || fw <= 0 || fh <= 0 {
			return nil
		}
		ox, _ := strconv.Atoi(k.fields[konfFieldOffsetX])
		oy, _ := strconv.Atoi(k.fields[konfFieldOffsetY])

		k.mriezka.frameWidth = fw
		k.mriezka.frameHeight = fh
		k.mriezka.offsetX = ox
		k.mriezka.offsetY = oy

		d := mriezkaJSON{FrameWidth: fw, FrameHeight: fh, OffsetX: ox, OffsetY: oy}
		if b, err := json.Marshal(d); err == nil {
			os.WriteFile(sidecarPath(k.mriezka.imagePath, k.mriezka.name), b, 0644)
		}
		return ebiten.Termination
	}

	return nil
}

func (k *mriezkaKonfig) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})

	if k.screenWidth == 0 || k.screenHeight == 0 {
		return
	}

	// Create ebiten image lazily within this RunGame session.
	if k.img == nil && k.rawImg != nil {
		k.img = ebiten.NewImageFromImage(k.rawImg)
	}

	const topPad = 25
	const bottomPad = 100

	imgW := k.mriezka.imageWidth
	imgH := k.mriezka.imageHeight
	availW := float64(k.screenWidth)
	availH := float64(k.screenHeight - topPad - bottomPad)

	scaleX := availW / float64(imgW)
	scaleY := availH / float64(imgH)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(imgW) * scale
	drawH := float64(imgH) * scale
	drawX := (float64(k.screenWidth) - drawW) / 2
	drawY := float64(topPad)

	if k.img != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(scale, scale)
		opts.GeoM.Translate(drawX, drawY)
		screen.DrawImage(k.img, opts)
	}

	// Grid overlay
	fw, _ := strconv.Atoi(k.fields[konfFieldSirka])
	fh, _ := strconv.Atoi(k.fields[konfFieldVyska])
	ox, _ := strconv.Atoi(k.fields[konfFieldOffsetX])
	oy, _ := strconv.Atoi(k.fields[konfFieldOffsetY])

	gridColor := color.RGBA{255, 220, 0, 200}
	if fw > 0 {
		for x := ox; x <= imgW; x += fw {
			sx := float32(drawX + float64(x)*scale)
			vector.StrokeLine(screen, sx, float32(drawY), sx, float32(drawY+drawH), 1, gridColor, false)
		}
	}
	if fh > 0 {
		for y := oy; y <= imgH; y += fh {
			sy := float32(drawY + float64(y)*scale)
			vector.StrokeLine(screen, float32(drawX), sy, float32(drawX+drawW), sy, 1, gridColor, false)
		}
	}

	// Instructions
	text.Draw(screen, "Nastav velikost snimku a stiskni Enter | Tab = dalsi pole | Shift+Tab = predchozi", basicfont.Face7x13, 5, 18, color.White)

	// Input fields
	baseY := k.screenHeight - bottomPad + 15
	for i, label := range konfFieldLabels {
		val := k.fields[i]
		var line string
		if i == k.active {
			line = label + ": [" + val + "|]"
		} else {
			line = label + ": " + val
		}
		col := color.RGBA{200, 200, 200, 255}
		if i == k.active {
			col = color.RGBA{255, 220, 0, 255}
		}
		text.Draw(screen, line, basicfont.Face7x13, 5, baseY+i*18, col)
	}

	// Validation hint
	fw2, err1 := strconv.Atoi(k.fields[konfFieldSirka])
	fh2, err2 := strconv.Atoi(k.fields[konfFieldVyska])
	if (k.fields[konfFieldSirka] != "" && (err1 != nil || fw2 <= 0)) ||
		(k.fields[konfFieldVyska] != "" && (err2 != nil || fh2 <= 0)) {
		text.Draw(screen, "Sirka a vyska musi byt kladne cislo!", basicfont.Face7x13, 5, baseY+4*18, color.RGBA{255, 100, 100, 255})
	}
}

func (k *mriezkaKonfig) Layout(outsideWidth, outsideHeight int) (int, int) {
	k.screenWidth = outsideWidth
	k.screenHeight = outsideHeight
	return outsideWidth, outsideHeight
}

func sidecarPath(imagePath, name string) string {
	if name == "" {
		return imagePath + ".mriezka"
	}
	return imagePath + "." + name + ".mriezka"
}

func runGridConfig(m *Mriezka, rawImg image.Image) {
	k := &mriezkaKonfig{
		mriezka: m,
		rawImg:  rawImg,
		fields:  [4]string{"", "", "0", "0"},
	}
	title := "Nastav mrizku: " + m.imagePath
	if m.name != "" {
		title += " (" + m.name + ")"
	}
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle(title)
	if err := ebiten.RunGame(k); err != nil && err != ebiten.Termination {
		log.Error("Chyba konfigurace mrizky", "chyba", err)
		os.Exit(1)
	}
}

func nactiMriezku(cesta, name string) (*Mriezka, bool) {
	data, err := os.ReadFile(sidecarPath(cesta, name))
	if err != nil {
		return nil, false
	}
	var d mriezkaJSON
	if err := json.Unmarshal(data, &d); err != nil || d.FrameWidth <= 0 || d.FrameHeight <= 0 {
		return nil, false
	}
	imgBytes, err := readFile(cesta)
	if err != nil {
		return nil, false
	}
	rawImg, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, false
	}
	return &Mriezka{
		name:        name,
		imagePath:   cesta,
		imageWidth:  rawImg.Bounds().Dx(),
		imageHeight: rawImg.Bounds().Dy(),
		frameWidth:  d.FrameWidth,
		frameHeight: d.FrameHeight,
		offsetX:     d.OffsetX,
		offsetY:     d.OffsetY,
	}, true
}

// TODO doc comment
func openGridConfigurator(image_path string, name string) (*Mriezka, error) {
	// No sidecar — decode image for config UI (raw, not through ebiten cache).
	imgBytes, err := readFile(image_path)
	if err != nil {
		return nil, err
	}
	rawImg, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Error("Nepodařilo se dekódovat obrázek mřížky", "chyba", err)
		os.Exit(1)
	}
	m := &Mriezka{
		name:        name,
		imagePath:   image_path,
		imageWidth:  rawImg.Bounds().Dx(),
		imageHeight: rawImg.Bounds().Dy(),
	}
	runGridConfig(m, rawImg)
	return m, nil
}
