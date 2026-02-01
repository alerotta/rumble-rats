package game

func NewWorld(matchId MatchId, width, height float64)* World{

	return &World{
		Match: matchId,
		Width: width,
		Height: height,
		Players:make(map[PlayerId]*Player) ,
		Tick: 0,
		Over: false,
		Winner: "",
	}
}

func (w *World) AddPlayer(id PlayerId, spawn Vec2) *Player{
	p := NewPlayer(id,spawn)
	w.Players[id] = p
	return p
}

func (w *World) GetPlayer(id PlayerId) *Player {
	return w.Players[id]
}
