package components

import (
	"hecate-shell/internal/installer/anim"
	"hecate-shell/internal/installer/styles"

	"github.com/charmbracelet/lipgloss"
)

// Option represents a selectable option
type Option struct {
	Label       string
	Description string
	Value       string
}

// Selector is an animated option selector (Yes/No, package manager, etc.)
type Selector struct {
	options  []Option
	selected int
	springs  []*anim.Spring
	focused  bool
}

// NewSelector creates a new selector with options
func NewSelector(options []Option) *Selector {
	springs := make([]*anim.Spring, len(options))
	for i := range springs {
		springs[i] = anim.NewBouncySpring()
		springs[i].SetPos(0)
	}

	return &Selector{
		options:  options,
		selected: 0,
		springs:  springs,
		focused:  true,
	}
}

// NewYesNoSelector creates a Yes/No selector
func NewYesNoSelector(defaultYes bool) *Selector {
	sel := NewSelector([]Option{
		{Label: "Yes", Value: "yes"},
		{Label: "No", Value: "no"},
	})
	if !defaultYes {
		sel.selected = 1
	}
	return sel
}

// Focus sets the focus state
func (s *Selector) Focus() {
	s.focused = true
}

// Blur removes focus
func (s *Selector) Blur() {
	s.focused = false
}

// Update advances animations
func (s *Selector) Update() {
	for i, spring := range s.springs {
		if i == s.selected {
			spring.SetTarget(1.0) // Selected = scaled up
		} else {
			spring.SetTarget(0.0) // Not selected = normal
		}
		spring.Update()
	}
}

// Next moves to the next option
func (s *Selector) Next() {
	s.selected = (s.selected + 1) % len(s.options)
}

// Prev moves to the previous option
func (s *Selector) Prev() {
	s.selected--
	if s.selected < 0 {
		s.selected = len(s.options) - 1
	}
}

// Selected returns the currently selected option
func (s *Selector) Selected() Option {
	return s.options[s.selected]
}

// SelectedIndex returns the index of the selected option
func (s *Selector) SelectedIndex() int {
	return s.selected
}

// SetSelected sets the selected index
func (s *Selector) SetSelected(idx int) {
	if idx >= 0 && idx < len(s.options) {
		s.selected = idx
	}
}

// View renders the selector
func (s *Selector) View() string {
	var rendered []string

	for i, opt := range s.options {
		springVal := s.springs[i].Pos()

		// Interpolate style based on spring value
		var style lipgloss.Style
		if springVal > 0.5 {
			style = styles.ButtonFocused
		} else {
			style = styles.Button
		}

		// Add extra padding based on spring (subtle bounce effect)
		extraPad := int(springVal * 1)
		if extraPad > 0 {
			style = style.PaddingLeft(3 + extraPad).PaddingRight(3 + extraPad)
		}

		label := opt.Label
		if i == 0 && opt.Value == "yes" {
			label += " (Recommended)"
		}

		rendered = append(rendered, style.Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, rendered...)
}

// ViewVertical renders the selector vertically (for lists)
func (s *Selector) ViewVertical() string {
	var rendered []string

	for i, opt := range s.options {
		springVal := s.springs[i].Pos()

		// Cursor indicator
		cursor := "  "
		if i == s.selected {
			cursor = styles.AccentText.Render("▸ ")
		}

		// Style based on selection
		var labelStyle lipgloss.Style
		if springVal > 0.5 {
			labelStyle = styles.AccentText.Bold(true)
		} else {
			labelStyle = styles.NormalText
		}

		line := cursor + labelStyle.Render(opt.Label)

		// Add description if present
		if opt.Description != "" {
			line += "\n    " + styles.DimText.Render(opt.Description)
		}

		rendered = append(rendered, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rendered...)
}
