package game

type PlayerId string
type MatchId string

type Vec2 struct {
	X float64
	Y float64
}

type Buttons struct {
	Up bool
	Down bool
	Left bool
	Right bool

	Attack bool
	Dash bool
}

type Input struct {

	Player PlayerId
	Seq uint32
	Tick uint32

	Buttons Buttons
}

type Player struct {
	ID PlayerId

	//position
	Pos Vec2
	Vel Vec2
	Facing Vec2

	//healt
	HP int
	Alive bool

	//cooldowns
	AttackCD float64
	DaskCD float64

}

type World struct {

	Match MatchId

	//sizes
	Width float64
	Height float64

	Players map[PlayerId]*Player
	Tick uint32

	Over bool
	Winner PlayerId
}

// very similar to player, it is possible that player contains server-only 
// value that the user should not be able to see, it can be useful later.

type PlayerSnapshot struct {
	ID PlayerId
	Pos Vec2
	Vel Vec2
	Facing Vec2
	HP int
	Alive bool
}

type Snapshot struct {
	Match MatchId
	Tick uint32

	Players []PlayerSnapshot

	Over bool
	Winner PlayerId
}