# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HecateShell is a custom Wayland shell built with QuickShell/QML for Arch Linux running the Niri compositor. The project consists of two main components:

1. **QML Shell** (root directory): The actual shell UI with hot-reloadable configuration
2. **Go CLI** (`hecate-shell-src/`): Management binary for installation, theming, and wallpaper control

## Development Commands

### Building the CLI
```bash
./build.sh
# or manually:
cd hecate-shell-src && go build -o hecate .
```

### Running the Shell
```bash
# Start shell (requires installation first)
./hecate run

# Reload shell (kills existing instance and starts fresh)
./hecate run --reload

# Run in debug mode (foreground, no daemonize)
./hecate run --debug
./hecate run -rd  # reload + debug
```

### Installing
```bash
# First time setup
./hecate install

# Force reinstall
./hecate install --force
```

### Theme Management
```bash
# Reload theme from theme.json
./hecate theme reload

# Set wallpaper with theme generation
./hecate wallpaper /path/to/image.jpg --generate-theme

# Set wallpaper with custom transition
./hecate wallpaper /path/to/image.jpg --transition fade --duration 2
```

### Updating
```bash
# Check for and pull updates
./hecate update

# Force update
./hecate update --force
```

### Manual Theme Generation
```bash
# From existing theme.json
matugen json ~/.config/HecateShell/theme.json -t scheme-fidelity -m dark --continue-on-error

# From wallpaper
matugen image ~/wallpaper.png -t scheme-fidelity -m dark --continue-on-error

# From color
matugen color '#D5C4AA' -t scheme-fidelity -m dark --continue-on-error
```

Note: `--continue-on-error` is required because matugen templates for other apps may expect color names that don't exist in all schemes.

## Architecture

### QML Shell Structure

```
shell.qml                    # Entry point - imports Niri, loads Bar
├── bar/Bar.qml             # Main bar container (left/center/right sections)
├── components/
│   ├── Workspaces.qml      # 5 workspace indicators with sliding animation
│   ├── DateWidget.qml      # Live-updating date display
│   ├── MusicVisualizer.qml # Real-time audio visualizer using cava
│   ├── TimeWidget.qml      # Live-updating time display
│   └── SystemWidgets.qml   # Audio (PipeWire) + WiFi (ping-based)
├── config/
│   ├── Config.qml          # Singleton config with hot-reload (checks every 1s)
│   ├── ColorLoader.qml     # Process-based theme.json reader
│   ├── cava.conf           # Audio visualizer configuration
│   └── templates/          # Matugen templates for app theming
│       ├── hecate.json     # Shell theme template
│       ├── cava.ini        # Cava audio visualizer
│       ├── spicetify.ini   # Spotify theme
│       ├── discord.css     # Vencord/system24 theme
│       └── micro.micro     # Micro editor colorscheme
└── theme.json              # Material Design 3 color palette (hot-reloaded)
```

### Go CLI Structure

```
hecate-shell-src/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go            # Cobra root command setup
│   ├── install.go         # Clone repo to ~/.config/HecateShell
│   ├── run.go             # Launch QuickShell with --daemonize
│   ├── update.go          # Git pull + version checking
│   ├── theme.go           # Theme reload via matugen
│   └── wallpaper.go       # swww integration + theme generation
└── internal/
    ├── config/config.go   # Path helpers, version checking
    ├── matugen/matugen.go # Merged matugen config handling
    └── shell/shell.go     # QuickShell process management (pkill)
```

### Hot-Reload System

The shell uses a timer-based hot-reload mechanism (config/Config.qml):
- Every 1 second, reads `$HOME/.config/HecateShell/theme.json` via ColorLoader
- Compares content with last read
- If changed, parses JSON and updates color properties
- All QML components reference `Config.` properties, so changes propagate immediately
- Paths are dynamically constructed using `Quickshell.env("HOME")` to work on any system

Color transitions are animated (300ms) across all widgets.

### Theme/Color Flow

1. Source colors stored in `theme.json` (Material Design 3 palette)
2. Matugen generates colors from wallpaper/color/existing palette
3. CLI merges user's `~/.config/matugen/config.toml` with HecateShell templates
4. Templates output to:
   - `theme.json` (shell colors)
   - `~/.config/cava/config` (visualizer)
   - `~/.config/spicetify/Themes/text/color.ini` (spotify)
   - `~/.config/Vencord/themes/sys24.css` (discord)
   - `~/.config/micro/colorschemes/matugen.micro` (editor)
5. Config.qml hot-reloads `theme.json` every 1 second
6. JSON is parsed and Material Design 3 colors mapped to shell properties
7. All widgets update via property bindings

Color mappings (defined in Config.qml):
- `backgroundColor` → `theme.surface`
- `backgroundColorAlt` → `theme.surfaceContainer`
- `borderColor` → `theme.outline`
- `textColor` → `theme.surfaceText`
- `textColorDim` → `theme.surfaceVariantText`
- `accentColor` → `theme.primary` (workspace indicator, visualizer, active icons)

### Matugen Integration

The CLI creates a temporary merged config when running matugen:
1. Reads user's `~/.config/matugen/config.toml` (if exists)
2. Appends HecateShell template definitions
3. Runs matugen with `-c <temp_config>` flag
4. Cleans up temp file after completion

This allows users to keep their own matugen templates while HecateShell adds its own.

### Audio Visualizer (MusicVisualizer.qml)

Uses `cava` with custom configuration:
- Frequency range: 50 Hz - 10 kHz
- Bars arranged: lows (left) → mids (center) → highs (right)
- Bars grow upward from bottom
- 20 bars total with 2px spacing
- Configuration in `config/cava.conf`

### Audio Control (SystemWidgets.qml)

PipeWire integration via `wpctl`:
- Click: toggle mute (between 0% and previous volume)
- Scroll: adjust volume (clamped 0-100%)
- Icon changes when muted (volume = 0%)

### Network Status (SystemWidgets.qml)

Ping-based connectivity check:
- Pings `1.1.1.1` to determine connection status
- Icon updates based on connectivity

### Workspaces (Workspaces.qml)

5 workspace indicators with:
- Sliding position indicator for active workspace
- Active workspace text fades to darker color
- Niri IPC integration for workspace state

## Key Technical Details

### QuickShell/Niri Integration

- Uses Niri QML module (imported as `Niri 0.1`)
- Shell connects to Niri compositor on startup
- Workspace state synchronized via Niri IPC

### CLI Installation Flow

1. `hecate install` clones repo from GitHub to `~/.config/HecateShell`
2. Checks for `shell.qml` to verify installation
3. `hecate run` launches QuickShell with config directory path
4. QuickShell runs with `--no-duplicate --daemonize` flags

### Update/Version System

- Version stored in `version` file (plain text)
- Remote version fetched from GitHub raw content
- Update compares versions, runs `git pull` if needed
- User can force update with `--force` flag

## Dependencies

Runtime:
- quickshell-git (Qt 6.10+, Niri support)
- cava (audio visualizer)
- pipewire + wireplumber (audio control)
- matugen (color generation)
- niri (Wayland compositor)
- swww (wallpaper daemon)

Build (Go CLI):
- Go 1.21+
- github.com/spf13/cobra v1.10.2

## Common Patterns

### Adding New QML Components

1. Create component in appropriate directory (e.g., `components/MyWidget.qml`)
2. Add to `qmldir` if creating a new module
3. Import in parent component (e.g., `Bar.qml`)
4. Reference `Config.` properties for theming
5. Use `Behavior on color { ColorAnimation { duration: 300 } }` for smooth transitions

### Adding New CLI Commands

1. Create new file in `cmd/` (e.g., `cmd/mycommand.go`)
2. Define cobra.Command with Use/Short/Long/RunE
3. Add to rootCmd in `init()` function
4. Use `internal/config` for path helpers
5. Use `internal/shell` for QuickShell process management

### Adding New Matugen Templates

1. Create template file in `config/templates/`
2. Add template definition in `internal/matugen/matugen.go`
3. Template uses matugen's `{{ colors.X.default.hex }}` syntax

### Modifying Color Scheme

- Edit `theme.json` with new Material Design 3 colors
- Run `hecate theme reload` or manually invoke matugen
- Shell updates automatically within 1 second

## Known Quirks

- QuickShell may be built against different Qt version than system (rebuild if needed)
- Matugen requires `--continue-on-error` due to templates expecting colors that may not exist in all schemes
- Hot-reload uses 1-second timer polling (could be improved with inotify)
- Repository URL hardcoded in `internal/config/config.go` as `river1337/HecateShell`
- ColorLoader uses `cat` command via Process component (QML limitation)
