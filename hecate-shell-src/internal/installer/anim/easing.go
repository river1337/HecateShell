package anim

import "math"

// Easing represents an easing animation
type Easing struct {
	start    float64
	end      float64
	duration int     // frames
	frame    int
	easeFunc EaseFunc
	done     bool
}

// EaseFunc is an easing function that takes progress (0-1) and returns eased value (0-1)
type EaseFunc func(t float64) float64

// Common easing functions
var (
	// EaseLinear - no easing
	EaseLinear EaseFunc = func(t float64) float64 { return t }

	// EaseInOutCubic - smooth ease in and out
	EaseInOutCubic EaseFunc = func(t float64) float64 {
		if t < 0.5 {
			return 4 * t * t * t
		}
		return 1 - math.Pow(-2*t+2, 3)/2
	}

	// EaseOutCubic - smooth ease out (starts fast, ends slow)
	EaseOutCubic EaseFunc = func(t float64) float64 {
		return 1 - math.Pow(1-t, 3)
	}

	// EaseInCubic - ease in (starts slow, ends fast)
	EaseInCubic EaseFunc = func(t float64) float64 {
		return t * t * t
	}

	// EaseOutQuad - quadratic ease out
	EaseOutQuad EaseFunc = func(t float64) float64 {
		return 1 - (1-t)*(1-t)
	}

	// EaseInOutQuad - quadratic ease in and out
	EaseInOutQuad EaseFunc = func(t float64) float64 {
		if t < 0.5 {
			return 2 * t * t
		}
		return 1 - math.Pow(-2*t+2, 2)/2
	}

	// EaseOutExpo - exponential ease out (very smooth)
	EaseOutExpo EaseFunc = func(t float64) float64 {
		if t == 1 {
			return 1
		}
		return 1 - math.Pow(2, -10*t)
	}
)

// NewEasing creates a new easing animation
func NewEasing(start, end float64, durationFrames int, easeFunc EaseFunc) *Easing {
	if easeFunc == nil {
		easeFunc = EaseInOutCubic
	}
	return &Easing{
		start:    start,
		end:      end,
		duration: durationFrames,
		frame:    0,
		easeFunc: easeFunc,
		done:     false,
	}
}

// NewSmoothEasing creates a smooth ease-in-out animation
func NewSmoothEasing(start, end float64, durationFrames int) *Easing {
	return NewEasing(start, end, durationFrames, EaseInOutCubic)
}

// Update advances the animation by one frame
func (e *Easing) Update() float64 {
	if e.done {
		return e.end
	}

	e.frame++
	if e.frame >= e.duration {
		e.frame = e.duration
		e.done = true
	}

	return e.Pos()
}

// Pos returns the current position
func (e *Easing) Pos() float64 {
	if e.duration == 0 {
		return e.end
	}
	t := float64(e.frame) / float64(e.duration)
	eased := e.easeFunc(t)
	return e.start + (e.end-e.start)*eased
}

// Done returns true if the animation is complete
func (e *Easing) Done() bool {
	return e.done
}

// Reset restarts the animation
func (e *Easing) Reset() {
	e.frame = 0
	e.done = false
}

// SetTarget changes the end position (resets animation)
func (e *Easing) SetTarget(end float64) {
	e.start = e.Pos() // Start from current position
	e.end = end
	e.frame = 0
	e.done = false
}

// Skip completes the animation immediately
func (e *Easing) Skip() {
	e.frame = e.duration
	e.done = true
}

// Easing2D is a 2D easing animation
type Easing2D struct {
	X *Easing
	Y *Easing
}

// NewEasing2D creates a new 2D easing animation
func NewEasing2D(startX, startY, endX, endY float64, durationFrames int, easeFunc EaseFunc) *Easing2D {
	return &Easing2D{
		X: NewEasing(startX, endX, durationFrames, easeFunc),
		Y: NewEasing(startY, endY, durationFrames, easeFunc),
	}
}

// NewSmoothEasing2D creates a smooth 2D easing animation
func NewSmoothEasing2D(startX, startY, endX, endY float64, durationFrames int) *Easing2D {
	return NewEasing2D(startX, startY, endX, endY, durationFrames, EaseInOutCubic)
}

// Update advances both animations
func (e *Easing2D) Update() (float64, float64) {
	return e.X.Update(), e.Y.Update()
}

// Pos returns current X,Y position
func (e *Easing2D) Pos() (float64, float64) {
	return e.X.Pos(), e.Y.Pos()
}

// Done returns true if both animations are complete
func (e *Easing2D) Done() bool {
	return e.X.Done() && e.Y.Done()
}

// Reset restarts both animations
func (e *Easing2D) Reset() {
	e.X.Reset()
	e.Y.Reset()
}

// SetTarget changes both targets
func (e *Easing2D) SetTarget(x, y float64) {
	e.X.SetTarget(x)
	e.Y.SetTarget(y)
}

// Skip completes both animations
func (e *Easing2D) Skip() {
	e.X.Skip()
	e.Y.Skip()
}
