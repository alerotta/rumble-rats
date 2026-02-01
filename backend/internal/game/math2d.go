package game

import "math"

func Add(a,b Vec2) Vec2{
	return  Vec2{
		X: a.X +b.X,
		Y: a.Y +b.Y,
	}
}

func Sub(a,b Vec2) Vec2{
	return Vec2{
		X : a.X - b.X,
		Y : a.Y - b.Y,
	}
}

func Mul(v Vec2,s float64) Vec2 {
	return Vec2{
		X: v.X * s,
		Y: v.Y * s,
	}
} 

func Len(v Vec2) float64{
	return math.Hypot(v.X,v.Y)
}

func Normalize(v Vec2) Vec2{
	lenght := Len(v)
	if lenght == 0{
		return Vec2{X: 0,Y: 0}
	}
	return Vec2{
		X: v.X /lenght,
		Y: v.Y /lenght,
	}
}

func Distance(a,b Vec2)float64{
	return Len(Sub(a,b))
}

func Clamp (x, min, max float64)float64{
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func ClampVec2 (v Vec2 , min , max Vec2) Vec2{
	return Vec2{
		X: Clamp(v.X, min.X , max.X ),
		Y: Clamp(v.Y, min.Y , max.Y ),
	}
}

func Approach(current, target, delta float64) float64 {
	if current < target {
		return math.Min(current+delta, target)
	}
	return math.Max(current-delta, target)
}
