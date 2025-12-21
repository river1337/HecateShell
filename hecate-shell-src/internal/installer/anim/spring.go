package anim

import "github.com/charmbracelet/harmonica"

// Spring wraps harmonica's spring for easy use
type Spring struct {
	spring harmonica.Spring
	pos    float64
	vel    float64
	target float64
}

// NewSpring creates a new spring animator
// frequency: how fast it oscillates (higher = faster)
// damping: how quickly it settles (1.0 = no overshoot, <1.0 = bouncy)
func NewSpring(frequency, damping float64) *Spring {
	return &Spring{
		spring: harmonica.NewSpring(harmonica.FPS(60), frequency, damping),
		pos:    0,
		vel:    0,
		target: 0,
	}
}

// NewBouncySpring creates a spring with nice bounce
func NewBouncySpring() *Spring {
	return NewSpring(6.0, 0.5) // Bouncy!
}

// NewSmoothSpring creates a spring with smooth motion (no overshoot)
func NewSmoothSpring() *Spring {
	return NewSpring(5.0, 1.0) // Critically damped
}

// NewSlowSpring creates a slower, gentler spring
func NewSlowSpring() *Spring {
	return NewSpring(3.0, 0.8)
}

// SetTarget sets the target position
func (s *Spring) SetTarget(target float64) {
	s.target = target
}

// SetPos sets the current position directly (for initialization)
func (s *Spring) SetPos(pos float64) {
	s.pos = pos
	s.vel = 0
}

// Update advances the spring simulation by one frame
func (s *Spring) Update() float64 {
	s.pos, s.vel = s.spring.Update(s.pos, s.vel, s.target)
	return s.pos
}

// Pos returns the current position
func (s *Spring) Pos() float64 {
	return s.pos
}

// AtRest returns true if the spring has settled
func (s *Spring) AtRest() bool {
	const epsilon = 0.001
	return abs(s.pos-s.target) < epsilon && abs(s.vel) < epsilon
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Spring2D is a 2D spring for X,Y movement
type Spring2D struct {
	X *Spring
	Y *Spring
}

// NewSpring2D creates a new 2D spring
func NewSpring2D(frequency, damping float64) *Spring2D {
	return &Spring2D{
		X: NewSpring(frequency, damping),
		Y: NewSpring(frequency, damping),
	}
}

// NewBouncySpring2D creates a bouncy 2D spring
func NewBouncySpring2D() *Spring2D {
	return &Spring2D{
		X: NewBouncySpring(),
		Y: NewBouncySpring(),
	}
}

// SetTarget sets both X and Y targets
func (s *Spring2D) SetTarget(x, y float64) {
	s.X.SetTarget(x)
	s.Y.SetTarget(y)
}

// SetPos sets both X and Y positions
func (s *Spring2D) SetPos(x, y float64) {
	s.X.SetPos(x)
	s.Y.SetPos(y)
}

// Update advances both springs
func (s *Spring2D) Update() (float64, float64) {
	return s.X.Update(), s.Y.Update()
}

// Pos returns current X,Y position
func (s *Spring2D) Pos() (float64, float64) {
	return s.X.Pos(), s.Y.Pos()
}

// AtRest returns true if both springs have settled
func (s *Spring2D) AtRest() bool {
	return s.X.AtRest() && s.Y.AtRest()
}
