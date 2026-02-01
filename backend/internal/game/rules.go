package game

func StepWorld (w *World, dt float64, inputs []Input){

	if w.Over{
		return
	}

	for _, in := range inputs{
		p:= w.Players[in.Player]
		if p == nil || !p.Alive {
			continue
		}
		move:=  MoveVectorFromButtons(in.Buttons)

		if Len(move) > 0{
			p.Facing = Normalize(move)
		}
		ApplyMovement(p,move,dt)
	}

	for _, p := range w.Players{

		if p == nil || !p.Alive {
			continue
		}

		p.Pos = Add(p.Pos, Mul(p.Vel, dt))

		p.Pos = ClampVec2(
			p.Pos,
			Vec2{X: 0 , Y: 0},
			Vec2{X: w.Width, Y : w.Height},
		)
	}

	w.Tick++
}

func MoveVectorFromButtons(b Buttons) Vec2 {

	// be careful at coordinates system of the screen
	var x,y float64
	if b.Left{
		x -=1
	}
	if b.Right{
		x += 1
	}
	if b.Up {
		y -= 1 
	}
	if b.Down{
		y += 1
	}
	return Vec2 {X: x , Y: y}
}

func ApplyMovement (p *Player, move Vec2, dt float64){
		const (
		accel     = 40.0 
		maxSpeed  = 8.0  
		friction  = 20.0
		deadZone  = 0.0001
	)
	
	if Len(move) > deadZone{
		dir := Normalize(move)
		p.Vel = Add(p.Vel, Mul(dir, accel*dt))
	} else {
		p.Vel.X = Approach(p.Vel.X, 0, friction*dt)
		p.Vel.Y = Approach(p.Vel.Y, 0, friction*dt)
	}
	speed := Len(p.Vel)
	if speed > maxSpeed {
		p.Vel = Mul(Normalize(p.Vel), maxSpeed)
	}
}

