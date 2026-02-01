package game

const (
	defaultHP = 100
)

func NewPlayer (id PlayerId, spawn Vec2) *Player {
	return &Player{
		ID: id,
		Pos: spawn,
		Vel: Vec2{0,0},
		Facing: Vec2{1,0},
		HP: defaultHP,
		Alive: true,

		AttackCD: 0,
		DaskCD: 0,
	}
}

func (p *Player) ApplyDamage (dmg int){
	if !p.Alive{
		return
	}
	p.HP -= dmg
	if p.HP <= 0 {
		p.HP = 0
		p.Alive = false
	}
}