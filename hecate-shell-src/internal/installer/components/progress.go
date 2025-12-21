package components

import (
	"fmt"
	"strings"

	"hecate-shell/internal/installer/anim"
	"hecate-shell/internal/installer/styles"
)

// Progress is an animated progress bar
type Progress struct {
	current    float64 // 0.0 to 1.0
	target     float64
	spring     *anim.Spring
	width      int
	label      string
	showPct    bool
}

// NewProgress creates a new progress bar
func NewProgress(width int) *Progress {
	spring := anim.NewSmoothSpring()
	spring.SetPos(0)
	return &Progress{
		current: 0,
		target:  0,
		spring:  spring,
		width:   width,
		showPct: true,
	}
}

// SetProgress sets the target progress (0.0 to 1.0)
func (p *Progress) SetProgress(pct float64) {
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	p.target = pct
	p.spring.SetTarget(pct)
}

// SetLabel sets the label shown above the progress bar
func (p *Progress) SetLabel(label string) {
	p.label = label
}

// ShowPercentage enables/disables percentage display
func (p *Progress) ShowPercentage(show bool) {
	p.showPct = show
}

// Update advances the animation
func (p *Progress) Update() {
	p.current = p.spring.Update()
}

// View renders the progress bar
func (p *Progress) View() string {
	filled := int(p.current * float64(p.width))
	if filled > p.width {
		filled = p.width
	}
	if filled < 0 {
		filled = 0
	}

	// Build the bar
	bar := strings.Repeat("█", filled)
	bar += strings.Repeat("░", p.width-filled)

	// Style it
	styledBar := styles.ProgressFill.Render(bar[:filled]) +
		styles.ProgressBg.Render(bar[filled:])

	result := styledBar

	// Add percentage
	if p.showPct {
		pct := fmt.Sprintf(" %3d%%", int(p.current*100))
		result += styles.DimText.Render(pct)
	}

	// Add label above
	if p.label != "" {
		result = styles.NormalText.Render(p.label) + "\n" + result
	}

	return result
}

// Done returns true if progress has reached 100%
func (p *Progress) Done() bool {
	return p.current >= 0.99 && p.spring.AtRest()
}

// TaskProgress tracks multiple tasks with overall progress
type TaskProgress struct {
	tasks       []Task
	current     int
	progress    *Progress
	taskSpinner *Spinner
}

// Task represents a single installation task
type Task struct {
	Name   string
	Status TaskStatus
}

// TaskStatus represents the state of a task
type TaskStatus int

const (
	TaskPending TaskStatus = iota
	TaskRunning
	TaskSuccess
	TaskFailed
)

// NewTaskProgress creates a new task progress tracker
func NewTaskProgress(tasks []string, width int) *TaskProgress {
	taskList := make([]Task, len(tasks))
	for i, name := range tasks {
		taskList[i] = Task{Name: name, Status: TaskPending}
	}

	return &TaskProgress{
		tasks:       taskList,
		current:     0,
		progress:    NewProgress(width),
		taskSpinner: NewSpinner(),
	}
}

// StartTask marks a task as running
func (tp *TaskProgress) StartTask(idx int) {
	if idx >= 0 && idx < len(tp.tasks) {
		tp.current = idx
		tp.tasks[idx].Status = TaskRunning
		tp.updateProgress()
	}
}

// CompleteTask marks a task as done
func (tp *TaskProgress) CompleteTask(idx int, success bool) {
	if idx >= 0 && idx < len(tp.tasks) {
		if success {
			tp.tasks[idx].Status = TaskSuccess
		} else {
			tp.tasks[idx].Status = TaskFailed
		}
		tp.updateProgress()
	}
}

func (tp *TaskProgress) updateProgress() {
	completed := 0
	for _, t := range tp.tasks {
		if t.Status == TaskSuccess || t.Status == TaskFailed {
			completed++
		}
	}
	tp.progress.SetProgress(float64(completed) / float64(len(tp.tasks)))
}

// Update advances animations
func (tp *TaskProgress) Update() {
	tp.progress.Update()
	tp.taskSpinner.Update()
}

// View renders the task progress
func (tp *TaskProgress) View() string {
	var lines []string

	// Task list
	for i, task := range tp.tasks {
		var icon string
		var labelStyle = styles.DimText

		switch task.Status {
		case TaskPending:
			icon = styles.DimText.Render("○")
		case TaskRunning:
			icon = styles.AccentText.Render(tp.taskSpinner.View())
			labelStyle = styles.NormalText
		case TaskSuccess:
			icon = styles.SuccessText.Render("✓")
			labelStyle = styles.SuccessText
		case TaskFailed:
			icon = styles.ErrorText.Render("✗")
			labelStyle = styles.ErrorText
		}

		line := fmt.Sprintf("  %s %s", icon, labelStyle.Render(task.Name))
		lines = append(lines, line)

		// Current task indicator
		if i == tp.current && task.Status == TaskRunning {
			lines[len(lines)-1] = styles.AccentText.Render("▸") + lines[len(lines)-1][1:]
		}
	}

	// Progress bar at bottom
	lines = append(lines, "")
	lines = append(lines, tp.progress.View())

	return strings.Join(lines, "\n")
}

// AllDone returns true if all tasks are complete
func (tp *TaskProgress) AllDone() bool {
	for _, t := range tp.tasks {
		if t.Status == TaskPending || t.Status == TaskRunning {
			return false
		}
	}
	return true
}

// HasErrors returns true if any task failed
func (tp *TaskProgress) HasErrors() bool {
	for _, t := range tp.tasks {
		if t.Status == TaskFailed {
			return true
		}
	}
	return false
}
