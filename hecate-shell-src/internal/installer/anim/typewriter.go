package anim

import (
	"math/rand"
	"strings"
)

// Typewriter animates text appearing character by character
type Typewriter struct {
	text       string
	visible    int
	frameCount int
	speed      int     // frames per character
	variation  float64 // random variation (0.0-1.0)
	done       bool
	cursor     bool
	cursorOn   bool
	cursorRate int // frames per cursor blink
}

// NewTypewriter creates a new typewriter effect
func NewTypewriter(text string) *Typewriter {
	return &Typewriter{
		text:       text,
		visible:    0,
		speed:      3, // 3 frames per char at 60fps = ~20 chars/sec
		variation:  0.3,
		cursor:     true,
		cursorRate: 30, // blink every 0.5 sec
	}
}

// WithSpeed sets characters per second
func (t *Typewriter) WithSpeed(charsPerSec int) *Typewriter {
	if charsPerSec > 0 {
		t.speed = 60 / charsPerSec
		if t.speed < 1 {
			t.speed = 1
		}
	}
	return t
}

// WithVariation sets the random speed variation (0.0-1.0)
func (t *Typewriter) WithVariation(v float64) *Typewriter {
	t.variation = v
	return t
}

// WithCursor enables/disables the cursor
func (t *Typewriter) WithCursor(enabled bool) *Typewriter {
	t.cursor = enabled
	return t
}

// Reset restarts the animation
func (t *Typewriter) Reset() {
	t.visible = 0
	t.frameCount = 0
	t.done = false
}

// SetText changes the text and resets
func (t *Typewriter) SetText(text string) {
	t.text = text
	t.Reset()
}

// Update advances the animation by one frame
func (t *Typewriter) Update() {
	t.frameCount++

	// Cursor blink
	if t.frameCount%t.cursorRate == 0 {
		t.cursorOn = !t.cursorOn
	}

	if t.done {
		return
	}

	// Calculate effective speed with variation
	effectiveSpeed := t.speed
	if t.variation > 0 {
		// Add random variation
		delta := int(float64(t.speed) * t.variation * (rand.Float64()*2 - 1))
		effectiveSpeed += delta
		if effectiveSpeed < 1 {
			effectiveSpeed = 1
		}
	}

	// Advance text
	if t.frameCount%effectiveSpeed == 0 {
		t.visible++
		if t.visible >= len([]rune(t.text)) {
			t.visible = len([]rune(t.text))
			t.done = true
		}
	}
}

// View returns the visible portion of text
func (t *Typewriter) View() string {
	runes := []rune(t.text)
	if t.visible > len(runes) {
		t.visible = len(runes)
	}

	visible := string(runes[:t.visible])

	// Add cursor
	if t.cursor {
		if t.cursorOn || !t.done {
			visible += "█"
		} else {
			visible += " "
		}
	}

	return visible
}

// Done returns true if the animation is complete
func (t *Typewriter) Done() bool {
	return t.done
}

// Skip completes the animation immediately
func (t *Typewriter) Skip() {
	t.visible = len([]rune(t.text))
	t.done = true
}

// MultiTypewriter handles multiple lines appearing in sequence
type MultiTypewriter struct {
	lines    []*Typewriter
	current  int
	lineGap  int // frames between lines
	gapCount int
}

// NewMultiTypewriter creates a multi-line typewriter
func NewMultiTypewriter(lines []string) *MultiTypewriter {
	tw := make([]*Typewriter, len(lines))
	for i, line := range lines {
		tw[i] = NewTypewriter(line)
	}
	return &MultiTypewriter{
		lines:   tw,
		lineGap: 10, // ~0.17 sec gap between lines
	}
}

// Update advances the animation
func (m *MultiTypewriter) Update() {
	if m.current >= len(m.lines) {
		return
	}

	m.lines[m.current].Update()

	if m.lines[m.current].Done() {
		m.gapCount++
		if m.gapCount >= m.lineGap {
			m.current++
			m.gapCount = 0
		}
	}
}

// View returns all visible lines
func (m *MultiTypewriter) View() string {
	var result strings.Builder
	for i := 0; i <= m.current && i < len(m.lines); i++ {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(m.lines[i].View())
	}
	return result.String()
}

// Done returns true if all lines are complete
func (m *MultiTypewriter) Done() bool {
	return m.current >= len(m.lines)
}

// Skip completes all animations immediately
func (m *MultiTypewriter) Skip() {
	for _, tw := range m.lines {
		tw.Skip()
	}
	m.current = len(m.lines)
}
