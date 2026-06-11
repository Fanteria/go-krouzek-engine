package gke

type blockSnapshot struct {
	coords    Coords
	isPostava bool
	velocityY float64
	onGround  bool
	actions   []Akce
}

func (g *game) saveSnapshot() {
	g.snapshot = make([]blockSnapshot, len(g.blocks))
	for i, b := range g.blocks {
		blk := b.getBlock()
		snap := blockSnapshot{coords: blk.coords}
		switch v := b.(type) {
		case *HratelnaPostava:
			snap.isPostava = true
			snap.velocityY = v.Postava.velocityY
			snap.onGround = v.Postava.onGround
			snap.actions = append([]Akce(nil), v.Postava.actualActions...)
		case *enemy:
			snap.isPostava = true
			snap.velocityY = v.Postava.velocityY
			snap.onGround = v.Postava.onGround
			snap.actions = append([]Akce(nil), v.Postava.actualActions...)
		}
		g.snapshot[i] = snap
	}
}

func (g *game) restoreSnapshot() {
	for i, b := range g.blocks {
		if i >= len(g.snapshot) {
			break
		}
		snap := g.snapshot[i]
		switch v := b.(type) {
		case *StatickyBlok:
			v.Blok.coords = snap.coords
		case *AnimovanyBlok:
			v.Blok.coords = snap.coords
		case *HratelnaPostava:
			v.Postava.Blok.coords = snap.coords
			v.Postava.velocityY = snap.velocityY
			v.Postava.onGround = snap.onGround
			v.Postava.actualActions = append([]Akce(nil), snap.actions...)
		case *enemy:
			v.Postava.Blok.coords = snap.coords
			v.Postava.velocityY = snap.velocityY
			v.Postava.onGround = snap.onGround
			v.Postava.actualActions = append([]Akce(nil), snap.actions...)
		}
	}
	g.camera.offsetX = 0
	g.camera.offsetY = 0
	g.animationIndex = 0
}
