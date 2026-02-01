package game

import (
	"context"
	"sync"
	"time"
)

type Engine struct {
	world *World
	tickRate int
	dt float64
	period time.Duration

	mu sync.Mutex
	pending []Input
}

func NewEngine (world *World, tickRate int) *Engine {
	if tickRate < 0 {
		tickRate = 60
	}
	return &Engine {
		world: world,
		tickRate: tickRate ,
		dt: 1.0/float64(tickRate),
		period: time.Second / time.Duration(tickRate),
		pending: make([]Input,0,256),

	}
}

func (e *Engine) World() *World {return e.world}
func (e *Engine) TickRate() int {return e.tickRate}
func (e *Engine) DT() float64 {return e.dt}

func (e* Engine) EnqueueInput (in Input){
	e.mu.Lock()
	e.pending = append(e.pending, in)
	e.mu.Unlock()
}

func (e* Engine) StepOnce() {
	e.mu.Lock()
	inputs := make([]Input, len(e.pending))
	copy(inputs, e.pending)
	e.pending = e.pending[:0]
	e.mu.Unlock()

	StepWorld(e.world, e.dt, inputs)
}

func (e* Engine) Run(ctx context.Context){
	ticker := time.NewTicker(e.period)
	defer ticker.Stop()

		for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.StepOnce()
		}
	}
}