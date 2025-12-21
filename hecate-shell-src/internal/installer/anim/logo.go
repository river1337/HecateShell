package anim

import (
	"math/rand"
	"strings"
)

// LogoChar represents a single character in the logo with its animation state
type LogoChar struct {
	char    rune
	finalX  int
	finalY  int
	easing  *Easing2D
	delay   int // frames to wait before starting
	started bool
}

// LogoAssembler animates characters flying in to form the logo
type LogoAssembler struct {
	chars      []LogoChar
	width      int
	height     int
	termWidth  int
	termHeight int
	frame      int
	done       bool
}

// NewLogoAssembler creates a new logo assembler from ASCII art
func NewLogoAssembler(logo string, termWidth, termHeight int) *LogoAssembler {
	lines := strings.Split(logo, "\n")

	// Remove empty first line if present
	if len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}

	// Find dimensions
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len([]rune(line)) > width {
			width = len([]rune(line))
		}
	}

	// Calculate center offset
	offsetX := (termWidth - width) / 2
	offsetY := (termHeight - height) / 2

	// Create chars with their final positions
	var chars []LogoChar
	for y, line := range lines {
		runes := []rune(line)
		for x, ch := range runes {
			if ch != ' ' {
				finalX := x + offsetX
				finalY := y + offsetY

				// Start from scattered positions (not too far)
				startX := float64(finalX) + (rand.Float64()*40 - 20)
				startY := float64(finalY) + (rand.Float64()*20 - 10)

				// Animation duration: 30-45 frames (~0.5-0.75 sec at 60fps)
				duration := 30 + rand.Intn(15)

				// Stagger delay based on position (left-to-right, top-to-bottom)
				delay := (x + y*2) / 3

				easing := NewEasing2D(startX, startY, float64(finalX), float64(finalY), duration, EaseOutCubic)

				chars = append(chars, LogoChar{
					char:   ch,
					finalX: finalX,
					finalY: finalY,
					easing: easing,
					delay:  delay,
				})
			}
		}
	}

	return &LogoAssembler{
		chars:      chars,
		width:      width,
		height:     height,
		termWidth:  termWidth,
		termHeight: termHeight,
	}
}

// Update advances all character animations
func (l *LogoAssembler) Update() {
	l.frame++

	allDone := true
	for i := range l.chars {
		// Check if this char should start animating
		if !l.chars[i].started && l.frame >= l.chars[i].delay {
			l.chars[i].started = true
		}

		// Update if started
		if l.chars[i].started {
			l.chars[i].easing.Update()
			if !l.chars[i].easing.Done() {
				allDone = false
			}
		} else {
			allDone = false
		}
	}
	l.done = allDone
}

// View renders the current state to a string
func (l *LogoAssembler) View() string {
	// Create a 2D buffer
	buffer := make([][]rune, l.termHeight)
	for i := range buffer {
		buffer[i] = make([]rune, l.termWidth)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	// Place each character at its current position
	for _, ch := range l.chars {
		if !ch.started {
			continue // Don't show chars that haven't started yet
		}

		x, y := ch.easing.Pos()
		ix, iy := int(x+0.5), int(y+0.5)

		// Bounds check
		if ix >= 0 && ix < l.termWidth && iy >= 0 && iy < l.termHeight {
			buffer[iy][ix] = ch.char
		}
	}

	// Convert buffer to string
	var result strings.Builder
	for _, row := range buffer {
		result.WriteString(string(row))
		result.WriteRune('\n')
	}
	return result.String()
}

// Done returns true when all characters have settled
func (l *LogoAssembler) Done() bool {
	return l.done
}

// Skip completes the animation immediately
func (l *LogoAssembler) Skip() {
	for i := range l.chars {
		l.chars[i].started = true
		l.chars[i].easing.Skip()
	}
	l.done = true
}

// Resize updates the terminal dimensions and recalculates positions
func (l *LogoAssembler) Resize(width, height int) {
	// Calculate new offset
	oldOffsetX := (l.termWidth - l.width) / 2
	oldOffsetY := (l.termHeight - l.height) / 2
	newOffsetX := (width - l.width) / 2
	newOffsetY := (height - l.height) / 2

	deltaX := float64(newOffsetX - oldOffsetX)
	deltaY := float64(newOffsetY - oldOffsetY)

	for i := range l.chars {
		l.chars[i].finalX += int(deltaX)
		l.chars[i].finalY += int(deltaY)
		l.chars[i].easing.SetTarget(float64(l.chars[i].finalX), float64(l.chars[i].finalY))
	}

	l.termWidth = width
	l.termHeight = height
}
