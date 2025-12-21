package components

import "hecate-shell/internal/installer/styles"

// Spinner is an animated loading spinner
type Spinner struct {
	frames     []string
	current    int
	frameRate  int // frames per spinner frame
	frameCount int
}

// NewSpinner creates a new spinner with default frames
func NewSpinner() *Spinner {
	return &Spinner{
		frames:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		frameRate: 6, // ~10 spinner frames per second at 60fps
	}
}

// NewDotsSpinner creates a dots spinner
func NewDotsSpinner() *Spinner {
	return &Spinner{
		frames:    []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		frameRate: 6,
	}
}

// NewLineSpinner creates a simple line spinner
func NewLineSpinner() *Spinner {
	return &Spinner{
		frames:    []string{"|", "/", "-", "\\"},
		frameRate: 8,
	}
}

// Update advances the spinner
func (s *Spinner) Update() {
	s.frameCount++
	if s.frameCount%s.frameRate == 0 {
		s.current = (s.current + 1) % len(s.frames)
	}
}

// View returns the current spinner frame
func (s *Spinner) View() string {
	return styles.AccentText.Render(s.frames[s.current])
}

// Reset resets the spinner to the first frame
func (s *Spinner) Reset() {
	s.current = 0
	s.frameCount = 0
}
