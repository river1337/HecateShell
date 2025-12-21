package components

import (
	"hecate-shell/internal/installer/anim"
	"hecate-shell/internal/installer/styles"

	"github.com/charmbracelet/lipgloss"
)

// CheckItem represents a checkable item
type CheckItem struct {
	Label       string
	Description string
	Checked     bool
	spring      *anim.Spring
}

// Checklist is an animated multi-select checklist
type Checklist struct {
	items       []CheckItem
	cursor      int
	cascadeIdx  int        // For cascade-in animation
	cascadeRate int        // Frames between each item appearing
	frameCount  int
	focused     bool
}

// NewChecklist creates a new checklist
func NewChecklist(items []CheckItem) *Checklist {
	// Add springs to each item
	for i := range items {
		items[i].spring = anim.NewBouncySpring()
		items[i].spring.SetPos(0)
		if items[i].Checked {
			items[i].spring.SetTarget(1)
		}
	}

	return &Checklist{
		items:       items,
		cursor:      0,
		cascadeIdx:  0,
		cascadeRate: 6, // ~10 items per second at 60fps
		focused:     true,
	}
}

// Focus sets focus state
func (c *Checklist) Focus() {
	c.focused = true
}

// Blur removes focus
func (c *Checklist) Blur() {
	c.focused = false
}

// Update advances animations
func (c *Checklist) Update() {
	c.frameCount++

	// Cascade animation
	if c.cascadeIdx < len(c.items) && c.frameCount%c.cascadeRate == 0 {
		c.cascadeIdx++
	}

	// Update springs
	for i := range c.items {
		if c.items[i].Checked {
			c.items[i].spring.SetTarget(1)
		} else {
			c.items[i].spring.SetTarget(0)
		}
		c.items[i].spring.Update()
	}
}

// Next moves cursor down
func (c *Checklist) Next() {
	c.cursor = (c.cursor + 1) % len(c.items)
}

// Prev moves cursor up
func (c *Checklist) Prev() {
	c.cursor--
	if c.cursor < 0 {
		c.cursor = len(c.items) - 1
	}
}

// Toggle toggles the current item
func (c *Checklist) Toggle() {
	c.items[c.cursor].Checked = !c.items[c.cursor].Checked
}

// ToggleAll toggles all items to the same state
func (c *Checklist) ToggleAll(checked bool) {
	for i := range c.items {
		c.items[i].Checked = checked
	}
}

// GetSelected returns labels of all checked items
func (c *Checklist) GetSelected() []string {
	var selected []string
	for _, item := range c.items {
		if item.Checked {
			selected = append(selected, item.Label)
		}
	}
	return selected
}

// Items returns all items
func (c *Checklist) Items() []CheckItem {
	return c.items
}

// View renders the checklist
func (c *Checklist) View() string {
	var lines []string

	for i, item := range c.items {
		// Only show items that have cascaded in
		if i >= c.cascadeIdx {
			continue
		}

		springVal := item.spring.Pos()

		// Cursor
		cursor := "  "
		if i == c.cursor && c.focused {
			cursor = styles.AccentText.Render("▸ ")
		}

		// Checkbox - interpolate between unchecked and checked
		var checkbox string
		if springVal > 0.5 {
			checkbox = styles.CheckboxChecked.Render("[✓]")
		} else {
			checkbox = styles.Checkbox.Render("[ ]")
		}

		// Label style based on checked state
		var labelStyle lipgloss.Style
		if item.Checked {
			labelStyle = styles.NormalText
		} else {
			labelStyle = styles.DimText
		}

		line := cursor + checkbox + " " + labelStyle.Render(item.Label)

		// Description
		if item.Description != "" {
			line += "\n      " + styles.DimText.Render(item.Description)
		}

		lines = append(lines, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// CascadeComplete returns true when all items have appeared
func (c *Checklist) CascadeComplete() bool {
	return c.cascadeIdx >= len(c.items)
}

// SkipCascade shows all items immediately
func (c *Checklist) SkipCascade() {
	c.cascadeIdx = len(c.items)
}
