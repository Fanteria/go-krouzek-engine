package gke

type camera struct {
	active      bool
	character   *Postava
	offsetX     float64
	offsetY     float64
	marginLeft  float64
	marginRight float64
	marginUp    float64
	marginDown  float64
}

func (k *camera) actualize(screenW, screenH int) {
	if !k.active || k.character == nil {
		k.offsetX = 0
		k.offsetY = 0
		return
	}
	p := k.character
	sub := p.getSubImage(game_instance.animationIndex)
	charW := float64(sub.Dx()) * p.scale.width
	charH := float64(sub.Dy()) * p.scale.height

	// X: move camera only when character crosses the margins
	charScreenX := p.coords.x - k.offsetX
	if charScreenX < k.marginLeft {
		k.offsetX = p.coords.x - k.marginLeft
	} else if charScreenX+charW > float64(screenW)-k.marginRight {
		k.offsetX = p.coords.x + charW - float64(screenW) + k.marginRight
	}
	if k.offsetX < 0 {
		k.offsetX = 0
	}

	// Y: move camera only when character crosses the margins
	charScreenY := p.coords.y - k.offsetY
	if charScreenY < k.marginUp {
		k.offsetY = p.coords.y - k.marginUp
	} else if charScreenY+charH > float64(screenH)-k.marginDown {
		k.offsetY = p.coords.y + charH - float64(screenH) + k.marginDown
	}
	if k.offsetY < 0 {
		k.offsetY = 0
	}
}
