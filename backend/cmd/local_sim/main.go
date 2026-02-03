package main

import (
	"fmt"
	"time"
	"github.com/alerotta/rumble-rats/backend/internal/game"
)

func main () {

	w:= game.NewWorld(game.MatchId("local-1"), 20,20)
	w.AddPlayer(game.PlayerId("player-1"),game.Vec2{X: 5,Y: 5})
	w.AddPlayer(game.PlayerId("player-2"),game.Vec2{X: 15,Y: 5})

	engine := game.NewEngine(w,60)
	ticker := time.NewTicker(time.Second/60)
	defer ticker.Stop()

	var seq uint32

	for range ticker.C {
		seq++

			engine.EnqueueInput(game.Input{
			Player: game.PlayerId("player-1"),
			Seq:    seq,
			Buttons: game.Buttons{
				Right: true,
			},
		})

				engine.EnqueueInput(game.Input{
			Player: game.PlayerId("player-2"),
			Seq:    seq,
			Buttons: game.Buttons{
				Left: true,
			},
		})

		engine.StepOnce()

		// 7) Observe state (print snapshot every 10 ticks)
		if engine.World().Tick%10 == 0 {
			snap := game.BuildSnapshot(engine.World())
			printSnapshot(snap)
		}

		// Optional: stop after some ticks
		if engine.World().Tick >= 300 {
			fmt.Println("done")
			return
		}



	}

}

func printSnapshot(s game.Snapshot) {
	fmt.Printf("tick=%d over=%v winner=%s\n", s.Tick, s.Over, s.Winner)
	for _, p := range s.Players {
		fmt.Printf("  %s pos=(%.2f,%.2f) vel=(%.2f,%.2f)\n",
			p.ID, p.Pos.X, p.Pos.Y, p.Vel.X, p.Vel.Y)
	}
}
