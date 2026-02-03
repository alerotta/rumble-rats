package game

import "sort"

func BuildSnapshot (w *World) Snapshot{

	s := Snapshot{
		Match: w.Match,
		Tick: w.Tick,
		Over: w.Over,
		Winner: w.Winner,
	}

	players := make([]PlayerSnapshot, 0 ,len(w.Players))
	for _ , p := range w.Players{
		if p == nil {
			continue
		}
		players = append(players, PlayerSnapshot{
			ID: p.ID,
			Pos: p.Pos,
			Vel: p.Vel,
			Facing: p.Facing,
			HP: p.HP,
			Alive: p.Alive,
		})
	}

	sort.Slice(players, func(i,j int)bool{
		return string(players[i].ID) <  string(players[j].ID)
	})

	s.Players = players
	return s
}