package tbc

type pool struct {
	casts []*Cast
}

func (p *pool) NewCast() *Cast {
	poolSize := len(p.casts)

	if poolSize <= 0 {
		p.fill()
		poolSize = len(p.casts)
	}

	c := p.casts[poolSize-1]
	p.casts = p.casts[:poolSize-1]
	return c
}

func (p *pool) fill() {
	newCasts := make([]Cast, 1000)
	for i := range newCasts {
		p.casts = append(p.casts, &newCasts[i])
	}
}

func (p *pool) ReturnCasts(casts []*Cast) {
	for _, v := range casts {
		v.Spell = nil
		v.IsLO = false
		v.IsClBounce = false
		v.CastTime = 0
		v.ManaCost = 0
		v.Hit = 0
		v.Crit = 0
		v.CritBonus = 0
		v.Spellpower = 0
		v.DidHit = false
		v.DidCrit = false
		v.DidDmg = 0
		v.CastAt = 0
		v.Effects = nil
	}

	p.casts = append(p.casts, casts...)
}
