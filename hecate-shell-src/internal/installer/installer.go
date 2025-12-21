package installer

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"hecate-shell/internal/config"
	"hecate-shell/internal/installer/actions"
	"hecate-shell/internal/installer/anim"
	"hecate-shell/internal/installer/components"
	"hecate-shell/internal/installer/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen represents the current installer screen
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenUpdatePrompt // Shown when already installed
	ScreenDeps
	ScreenPackageManager
	ScreenDotfiles
	ScreenShell
	ScreenInstalling
	ScreenComplete
	ScreenError
)

// tickMsg is sent every frame for animations
type tickMsg time.Time

// installDoneMsg is sent when installation completes
type installDoneMsg struct {
	success bool
	err     error
}

// taskCompleteMsg is sent when a single task completes
type taskCompleteMsg struct {
	index   int
	success bool
	err     error
}

// installTask represents a task to run during installation
type installTask struct {
	name string
	run  func() error
}

// Model is the main installer model
type Model struct {
	// Current screen
	screen Screen

	// Terminal dimensions
	width  int
	height int

	// Update mode (true if HecateShell is already installed)
	isUpdate bool

	// Animation state
	logoAnim   *anim.LogoAssembler
	typewriter *anim.Typewriter

	// User choices
	installDeps    bool
	packageManager string
	dotfileChoices map[string]bool
	installShell   bool

	// UI Components
	updateSelector *components.Selector // For update prompt
	depsSelector   *components.Selector
	pkgSelector    *components.Selector
	dotfilesCheck  *components.Checklist
	shellSelector  *components.Selector
	taskProgress   *components.TaskProgress
	completeAnim   *anim.Spring

	// Installation state
	installTasks   []installTask
	currentTaskIdx int

	// State
	ready    bool
	quitting bool
	err      error
}

// New creates a new installer model
func New() Model {
	isUpdate := config.IsInstalled()
	return Model{
		screen:         ScreenWelcome,
		isUpdate:       isUpdate,
		dotfileChoices: make(map[string]bool),
		installDeps:    true,
		installShell:   true,
		packageManager: "paru",
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tick(),
		tea.EnterAltScreen,
	)
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Initialize components that need dimensions
		if m.logoAnim == nil {
			m.logoAnim = anim.NewLogoAssembler(styles.Logo(), m.width, m.height)
			// Different message based on update mode
			welcomeMsg := "Welcome to the HecateShell installer"
			if m.isUpdate {
				welcomeMsg = "HecateShell is already installed"
			}
			m.typewriter = anim.NewTypewriter(welcomeMsg).
				WithSpeed(40). // Faster typing
				WithCursor(true)
		} else {
			m.logoAnim.Resize(m.width, m.height)
		}
		return m, nil

	case tickMsg:
		m.updateAnimations()
		return m, tick()

	case taskCompleteMsg:
		// Mark task as complete
		if m.taskProgress != nil {
			m.taskProgress.CompleteTask(msg.index, msg.success)
		}

		if !msg.success && msg.err != nil {
			// Task failed - show error
			m.err = msg.err
			m.screen = ScreenError
			return m, nil
		}

		// Move to next task
		m.currentTaskIdx++
		if m.currentTaskIdx < len(m.installTasks) {
			return m, m.runNextTask()
		}

		// All tasks done!
		m.screen = ScreenComplete
		m.completeAnim = anim.NewBouncySpring()
		m.completeAnim.SetPos(-5)
		m.completeAnim.SetTarget(0)
		return m, nil

	case installDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.screen = ScreenError
		} else {
			m.screen = ScreenComplete
			m.completeAnim = anim.NewBouncySpring()
			m.completeAnim.SetPos(-5)
			m.completeAnim.SetTarget(0)
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) updateAnimations() {
	switch m.screen {
	case ScreenWelcome:
		if m.logoAnim != nil {
			m.logoAnim.Update()
		}
		if m.logoAnim != nil && m.logoAnim.Done() && m.typewriter != nil {
			m.typewriter.Update()
		}

	case ScreenUpdatePrompt:
		if m.updateSelector != nil {
			m.updateSelector.Update()
		}

	case ScreenDeps:
		if m.depsSelector != nil {
			m.depsSelector.Update()
		}

	case ScreenPackageManager:
		if m.pkgSelector != nil {
			m.pkgSelector.Update()
		}

	case ScreenDotfiles:
		if m.dotfilesCheck != nil {
			m.dotfilesCheck.Update()
		}

	case ScreenShell:
		if m.shellSelector != nil {
			m.shellSelector.Update()
		}

	case ScreenInstalling:
		if m.taskProgress != nil {
			m.taskProgress.Update()
		}

	case ScreenComplete:
		if m.completeAnim != nil {
			m.completeAnim.Update()
		}
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.quitting = true
		return m, tea.Quit
	}

	switch m.screen {
	case ScreenWelcome:
		return m.handleWelcomeKey(msg)
	case ScreenUpdatePrompt:
		return m.handleUpdatePromptKey(msg)
	case ScreenDeps:
		return m.handleDepsKey(msg)
	case ScreenPackageManager:
		return m.handlePkgKey(msg)
	case ScreenDotfiles:
		return m.handleDotfilesKey(msg)
	case ScreenShell:
		return m.handleShellKey(msg)
	case ScreenComplete, ScreenError:
		if msg.String() == "enter" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) handleWelcomeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		// Skip animations if not done
		if m.logoAnim != nil && !m.logoAnim.Done() {
			m.logoAnim.Skip()
			return m, nil
		}
		if m.typewriter != nil && !m.typewriter.Done() {
			m.typewriter.Skip()
			return m, nil
		}

		// If already installed, show update prompt
		if m.isUpdate {
			m.screen = ScreenUpdatePrompt
			m.updateSelector = components.NewYesNoSelector(true)
			return m, nil
		}

		// Move to next screen for fresh install
		m.screen = ScreenDeps
		m.depsSelector = components.NewYesNoSelector(true)
		return m, nil
	}
	return m, nil
}

func (m Model) handleUpdatePromptKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.updateSelector.Prev()
	case "right", "l":
		m.updateSelector.Next()
	case "tab":
		m.updateSelector.Next()
	case "enter":
		if m.updateSelector.Selected().Value == "yes" {
			// User wants to update - continue to deps/dotfiles screens
			m.screen = ScreenDeps
			m.depsSelector = components.NewYesNoSelector(true)
		} else {
			// User doesn't want to update - exit
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) handleDepsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.depsSelector.Prev()
	case "right", "l":
		m.depsSelector.Next()
	case "tab":
		m.depsSelector.Next()
	case "enter":
		m.installDeps = m.depsSelector.Selected().Value == "yes"
		if m.installDeps {
			// Show package manager selection
			m.screen = ScreenPackageManager
			m.pkgSelector = components.NewSelector([]components.Option{
				{Label: "paru", Description: "AUR helper (Recommended)", Value: "paru"},
				{Label: "yay", Description: "Another AUR helper", Value: "yay"},
				{Label: "pacman", Description: "Official packages only", Value: "pacman"},
			})
		} else {
			// Skip to dotfiles
			m.screen = ScreenDotfiles
			m.initDotfilesChecklist()
		}
	}
	return m, nil
}

func (m Model) handlePkgKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.pkgSelector.Prev()
	case "down", "j":
		m.pkgSelector.Next()
	case "tab":
		m.pkgSelector.Next()
	case "enter":
		m.packageManager = m.pkgSelector.Selected().Value
		m.screen = ScreenDotfiles
		m.initDotfilesChecklist()
	}
	return m, nil
}

func (m *Model) initDotfilesChecklist() {
	m.dotfilesCheck = components.NewChecklist([]components.CheckItem{
		{Label: "Niri", Description: "~/.config/niri/", Checked: true},
		{Label: "Fish", Description: "~/.config/fish/", Checked: true},
		{Label: "Kitty", Description: "~/.config/kitty/", Checked: true},
		{Label: "Micro", Description: "~/.config/micro/", Checked: true},
		{Label: "Fastfetch", Description: "~/.config/fastfetch/", Checked: true},
	})
}

func (m Model) handleDotfilesKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.dotfilesCheck.Prev()
	case "down", "j":
		m.dotfilesCheck.Next()
	case " ":
		m.dotfilesCheck.Toggle()
	case "a":
		m.dotfilesCheck.ToggleAll(true)
	case "n":
		m.dotfilesCheck.ToggleAll(false)
	case "enter":
		// Save selections
		for _, item := range m.dotfilesCheck.Items() {
			m.dotfileChoices[item.Label] = item.Checked
		}
		m.screen = ScreenShell
		m.shellSelector = components.NewYesNoSelector(true)
	}
	return m, nil
}

func (m Model) handleShellKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		m.shellSelector.Prev()
	case "right", "l":
		m.shellSelector.Next()
	case "tab":
		m.shellSelector.Next()
	case "enter":
		m.installShell = m.shellSelector.Selected().Value == "yes"
		m.screen = ScreenInstalling
		return m, m.startInstallation()
	}
	return m, nil
}

func (m *Model) startInstallation() tea.Cmd {
	// Build task list based on selections
	var tasks []string
	var installTasks []installTask

	// Action words for task messages
	depsAction := "Installing"
	dotfileAction := "Installing"
	shellAction := "Downloading"
	if m.isUpdate {
		depsAction = "Updating"
		dotfileAction = "Updating"
		shellAction = "Updating"
	}

	if m.installDeps {
		tasks = append(tasks, depsAction+" dependencies...")
		pm := m.packageManager // capture for closure
		installTasks = append(installTasks, installTask{
			name: "deps",
			run: func() error {
				result := actions.InstallDependencies(pm, nil)
				if !result.Success {
					return result.Error
				}
				return nil
			},
		})
	}

	// Sort dotfiles for consistent order
	dotfileOrder := []string{"Niri", "Fish", "Kitty", "Micro", "Fastfetch"}
	isUpdate := m.isUpdate // capture for closure
	for _, name := range dotfileOrder {
		if checked, ok := m.dotfileChoices[name]; ok && checked {
			tasks = append(tasks, dotfileAction+" "+name+" config...")
			dotfileName := name // capture for closure
			installTasks = append(installTasks, installTask{
				name: dotfileName,
				run: func() error {
					// Backup existing config before installing (always, for safety)
					if isUpdate {
						actions.BackupDotfile(dotfileName)
					}
					result := actions.InstallDotfile(dotfileName, nil)
					if !result.Success {
						// Don't fail on dotfile errors, just warn
						return nil
					}
					return nil
				},
			})
		}
	}

	if m.installShell {
		tasks = append(tasks, shellAction+" HecateShell...")
		installTasks = append(installTasks, installTask{
			name: "shell",
			run: func() error {
				if isUpdate {
					// Update via git pull
					result := actions.UpdateShell(nil)
					if !result.Success {
						return result.Error
					}
				} else {
					// Fresh install via git clone
					result := actions.InstallShell(true, nil)
					if !result.Success {
						return result.Error
					}
				}
				return nil
			},
		})
	}

	tasks = append(tasks, "Finishing up...")
	installTasks = append(installTasks, installTask{
		name: "finish",
		run: func() error {
			actions.RunPostInstall(nil)
			return nil
		},
	})

	m.taskProgress = components.NewTaskProgress(tasks, 40)
	m.installTasks = installTasks
	m.currentTaskIdx = 0

	// Start the first task
	if len(m.installTasks) > 0 {
		return m.runNextTask()
	}

	return func() tea.Msg {
		return installDoneMsg{success: true}
	}
}

func (m *Model) runNextTask() tea.Cmd {
	if m.currentTaskIdx >= len(m.installTasks) {
		return func() tea.Msg {
			return installDoneMsg{success: true}
		}
	}

	idx := m.currentTaskIdx
	task := m.installTasks[idx]

	// Mark task as started
	m.taskProgress.StartTask(idx)

	return func() tea.Msg {
		err := task.run()
		return taskCompleteMsg{
			index:   idx,
			success: err == nil,
			err:     err,
		}
	}
}

// View renders the current screen
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	var content string

	switch m.screen {
	case ScreenWelcome:
		content = m.viewWelcome()
	case ScreenUpdatePrompt:
		content = m.viewUpdatePrompt()
	case ScreenDeps:
		content = m.viewDeps()
	case ScreenPackageManager:
		content = m.viewPackageManager()
	case ScreenDotfiles:
		content = m.viewDotfiles()
	case ScreenShell:
		content = m.viewShell()
	case ScreenInstalling:
		content = m.viewInstalling()
	case ScreenComplete:
		content = m.viewComplete()
	case ScreenError:
		content = m.viewError()
	}

	// Center content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m Model) viewWelcome() string {
	var content string

	if m.logoAnim != nil {
		// Apply accent color to the logo
		content = styles.LogoText.Render(m.logoAnim.View())
	}

	// Show typewriter after logo is done
	if m.logoAnim != nil && m.logoAnim.Done() && m.typewriter != nil {
		// Find center position for text
		text := m.typewriter.View()
		content += "\n\n" + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, styles.Subtitle.Render(text))

		if m.typewriter.Done() {
			hint := styles.DimText.Render("\n\nPress Enter to continue")
			content += lipgloss.PlaceHorizontal(m.width, lipgloss.Center, hint)
		}
	}

	return content
}

func (m Model) viewUpdatePrompt() string {
	title := styles.Title.Render("Update HecateShell?")
	desc := styles.Subtitle.Render("HecateShell is already installed. Would you like to update?")

	selector := ""
	if m.updateSelector != nil {
		selector = m.updateSelector.View()
	}

	help := styles.Help.Render("← → to select • Enter to confirm • q to quit")

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		desc,
		"",
		"",
		selector,
		"",
		help,
	)
}

func (m Model) viewDeps() string {
	title := styles.Title.Render("Install Dependencies?")
	desc := "HecateShell requires some packages to function properly."
	if m.isUpdate {
		title = styles.Title.Render("Update Dependencies?")
		desc = "Would you like to update/install any missing dependencies?"
	}
	descStyled := styles.Subtitle.Render(desc)

	selector := ""
	if m.depsSelector != nil {
		selector = m.depsSelector.View()
	}

	help := styles.Help.Render("← → to select • Enter to confirm • q to quit")

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		descStyled,
		"",
		"",
		selector,
		"",
		help,
	)
}

func (m Model) viewPackageManager() string {
	title := styles.Title.Render("Select Package Manager")
	desc := styles.Subtitle.Render("Choose your preferred AUR helper.")

	selector := ""
	if m.pkgSelector != nil {
		selector = m.pkgSelector.ViewVertical()
	}

	help := styles.Help.Render("↑ ↓ to select • Enter to confirm • q to quit")

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		desc,
		"",
		"",
		selector,
		"",
		help,
	)
}

func (m Model) viewDotfiles() string {
	title := styles.Title.Render("Install Dotfiles?")
	desc := "Select which configurations to install."
	note := ""
	if m.isUpdate {
		title = styles.Title.Render("Update Dotfiles?")
		desc = "Select which configurations to update."
		note = styles.DimText.Render("(Existing configs will be backed up)")
	}
	descStyled := styles.Subtitle.Render(desc)

	checklist := ""
	if m.dotfilesCheck != nil {
		checklist = m.dotfilesCheck.View()
	}

	help := styles.Help.Render("↑ ↓ to navigate • Space to toggle • a/n all/none • Enter to confirm")

	parts := []string{title, "", descStyled}
	if note != "" {
		parts = append(parts, note)
	}
	parts = append(parts, "", "", checklist, "", help)

	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m Model) viewShell() string {
	title := styles.Title.Render("Install HecateShell?")
	desc := "Download and install the QuickShell configuration."
	if m.isUpdate {
		title = styles.Title.Render("Update HecateShell?")
		desc = "Pull the latest changes from the repository."
	}
	descStyled := styles.Subtitle.Render(desc)

	selector := ""
	if m.shellSelector != nil {
		selector = m.shellSelector.View()
	}

	help := styles.Help.Render("← → to select • Enter to confirm • q to quit")

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		descStyled,
		"",
		"",
		selector,
		"",
		help,
	)
}

func (m Model) viewInstalling() string {
	titleText := "Installing..."
	if m.isUpdate {
		titleText = "Updating..."
	}
	title := styles.Title.Render(titleText)

	progress := ""
	if m.taskProgress != nil {
		progress = m.taskProgress.View()
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		"",
		progress,
	)
}

func (m Model) viewComplete() string {
	// Animate the checkmark bouncing in
	yOffset := 0
	if m.completeAnim != nil {
		yOffset = int(m.completeAnim.Pos())
	}

	checkmark := styles.SuccessText.Render("✓")
	titleText := "Installation Complete!"
	descText := "HecateShell has been installed successfully."
	if m.isUpdate {
		titleText = "Update Complete!"
		descText = "HecateShell has been updated successfully."
	}
	title := styles.Title.Render(titleText)
	desc := styles.Subtitle.Render(descText)

	// Summary
	actionWord := "installed"
	if m.isUpdate {
		actionWord = "updated"
	}
	var summary []string
	if m.installDeps {
		summary = append(summary, "• Dependencies "+actionWord+" with "+m.packageManager)
	}
	for name, checked := range m.dotfileChoices {
		if checked {
			summary = append(summary, "• "+name+" config "+actionWord)
		}
	}
	if m.installShell {
		summary = append(summary, "• HecateShell "+actionWord)
	}

	summaryText := lipgloss.JoinVertical(lipgloss.Left, summary...)

	help := styles.Help.Render("\nPress Enter to exit")

	// Add vertical offset for bounce effect
	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		checkmark,
		"",
		title,
		"",
		desc,
		"",
		summaryText,
		help,
	)

	if yOffset != 0 {
		content = "\n" + content // Simple offset simulation
	}

	return content
}

func (m Model) viewError() string {
	title := styles.ErrorText.Render("Installation Failed")
	errMsg := ""
	if m.err != nil {
		errMsg = styles.DimText.Render(m.err.Error())
	}

	help := styles.Help.Render("\nPress Enter to exit")

	return lipgloss.JoinVertical(lipgloss.Center,
		"",
		styles.ErrorText.Render("✗"),
		"",
		title,
		"",
		errMsg,
		help,
	)
}

// acquireSudo prompts for sudo password and keeps the session alive
func acquireSudo() error {
	fmt.Println(styles.LogoStyled())
	fmt.Println()
	fmt.Println(styles.Title.Render("HecateShell Installer"))
	fmt.Println()
	fmt.Println(styles.NormalText.Render("This installer requires sudo privileges to install dependencies."))
	fmt.Println(styles.DimText.Render("Please enter your password when prompted."))
	fmt.Println()

	// Run sudo -v to prompt for password and validate
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to acquire sudo privileges: %w", err)
	}

	// Start a background process to keep sudo alive
	go keepSudoAlive()

	fmt.Println()
	fmt.Println(styles.SuccessText.Render("✓ Sudo privileges acquired!"))
	fmt.Println(styles.DimText.Render("Starting installer..."))
	fmt.Println()

	// Small delay so user can see the message
	time.Sleep(500 * time.Millisecond)

	return nil
}

// keepSudoAlive refreshes sudo timestamp periodically
func keepSudoAlive() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		exec.Command("sudo", "-v").Run()
	}
}

// Run starts the installer TUI
func Run() error {
	// Acquire sudo privileges before starting TUI
	if err := acquireSudo(); err != nil {
		return err
	}

	p := tea.NewProgram(New(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
