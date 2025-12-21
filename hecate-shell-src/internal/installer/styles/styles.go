package styles

import "github.com/charmbracelet/lipgloss"

// Theme colors (Gruvbox Dark)
var (
	// Gruvbox palette
	Bg0Hard  = lipgloss.Color("#1d2021") // Darkest background
	Bg0      = lipgloss.Color("#282828") // Main background
	Bg1      = lipgloss.Color("#3c3836") // Lighter background
	Bg2      = lipgloss.Color("#504945") // Even lighter
	Fg0      = lipgloss.Color("#fbf1c7") // Brightest foreground
	Fg1      = lipgloss.Color("#ebdbb2") // Main foreground
	Fg2      = lipgloss.Color("#d5c4a1") // Dimmer foreground
	Fg3      = lipgloss.Color("#bdae93") // Even dimmer
	Fg4      = lipgloss.Color("#a89984") // Dimmest foreground
	Gray     = lipgloss.Color("#928374") // Gray

	// Gruvbox accent colors
	Red       = lipgloss.Color("#fb4934")
	Green     = lipgloss.Color("#b8bb26")
	Yellow    = lipgloss.Color("#fabd2f")
	Blue      = lipgloss.Color("#83a598")
	Purple    = lipgloss.Color("#d3869b")
	Aqua      = lipgloss.Color("#8ec07c")
	Orange    = lipgloss.Color("#fe8019")

	// Semantic mappings
	Primary    = Orange   // Main accent (Gruvbox orange)
	Secondary  = Aqua     // Secondary accent
	Surface    = Bg0      // Background
	SurfaceAlt = Bg1      // Elevated surface
	Text       = Fg1      // Main text
	TextDim    = Fg4      // Dimmed text
	Success    = Green    // Success state
	Warning    = Yellow   // Warning state
	Error      = Red      // Error state
)

// Base styles
var (
	// Container for full screen
	Container = lipgloss.NewStyle().
		Background(Surface)

	// Title text
	Title = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	// Subtitle/description
	Subtitle = lipgloss.NewStyle().
		Foreground(TextDim)

	// Normal text
	NormalText = lipgloss.NewStyle().
		Foreground(Text)

	// Dimmed text
	DimText = lipgloss.NewStyle().
		Foreground(TextDim)

	// Success text
	SuccessText = lipgloss.NewStyle().
		Foreground(Success)

	// Error text
	ErrorText = lipgloss.NewStyle().
		Foreground(Error)

	// Warning text
	WarningText = lipgloss.NewStyle().
		Foreground(Warning)

	// Accent text
	AccentText = lipgloss.NewStyle().
		Foreground(Primary)

	// Logo text (uses accent color)
	LogoText = lipgloss.NewStyle().
		Foreground(Primary)

	// Button (unfocused)
	Button = lipgloss.NewStyle().
		Foreground(Text).
		Background(SurfaceAlt).
		Padding(0, 3).
		MarginRight(2)

	// Button (focused)
	ButtonFocused = lipgloss.NewStyle().
		Foreground(Bg0).
		Background(Primary).
		Padding(0, 3).
		MarginRight(2).
		Bold(true)

	// Checkbox unchecked
	Checkbox = lipgloss.NewStyle().
		Foreground(TextDim)

	// Checkbox checked
	CheckboxChecked = lipgloss.NewStyle().
		Foreground(Success)

	// Progress bar background
	ProgressBg = lipgloss.NewStyle().
		Foreground(SurfaceAlt)

	// Progress bar fill
	ProgressFill = lipgloss.NewStyle().
		Foreground(Primary)

	// Card/box
	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SurfaceAlt).
		Padding(1, 2)

	// Card focused
	CardFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Primary).
		Padding(1, 2)

	// Help text at bottom
	Help = lipgloss.NewStyle().
		Foreground(TextDim).
		MarginTop(1)
)

// Logo returns the HecateShell ASCII logo
func Logo() string {
	return `
 ██╗  ██╗███████╗ ██████╗ █████╗ ████████╗███████╗
 ██║  ██║██╔════╝██╔════╝██╔══██╗╚══██╔══╝██╔════╝
 ███████║█████╗  ██║     ███████║   ██║   █████╗
 ██╔══██║██╔══╝  ██║     ██╔══██║   ██║   ██╔══╝
 ██║  ██║███████╗╚██████╗██║  ██║   ██║   ███████╗
 ╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝   ╚═╝   ╚══════╝
           ███████╗██╗  ██╗███████╗██╗     ██╗
           ██╔════╝██║  ██║██╔════╝██║     ██║
           ███████╗███████║█████╗  ██║     ██║
           ╚════██║██╔══██║██╔══╝  ██║     ██║
           ███████║██║  ██║███████╗███████╗███████╗
           ╚══════╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝`
}

// LogoStyled returns the logo with styling applied
func LogoStyled() string {
	return LogoText.Render(Logo())
}
