# HecateShell

A modular, customizable Wayland shell built with QuickShell for Arch Linux running Niri.

## Features

- **Hot-reloading themes** - Edit `theme.json` and see changes instantly
- **Material Design 3 colors** - Generated from wallpapers via matugen
- **Integrated theming** - Syncs colors to cava, spicetify, discord, and micro
- **Audio visualizer** - Real-time cava integration in the bar
- **PipeWire audio control** - Click to mute, scroll to adjust volume
- **Workspace indicators** - Animated workspace switching via Niri IPC

## Installation

```bash
# Build the CLI
./build.sh

# Install (clones to ~/.config/HecateShell)
./hecate install

# Start the shell
./hecate run
```

## Usage

```bash
# Start shell (daemonized)
hecate run

# Start in debug mode (foreground)
hecate run --debug

# Reload shell
hecate run --reload

# Set wallpaper
hecate wallpaper /path/to/image.jpg

# Set wallpaper and generate theme
hecate wallpaper /path/to/image.jpg --generate-theme

# Reload theme from theme.json
hecate theme reload

# Check for updates
hecate update
```

## Project Structure

```
.
├── shell.qml              # Main entry point
├── qmldir                 # QML module definitions
├── theme.json             # Material Design 3 color palette
├── config/
│   ├── Config.qml         # Singleton config with hot-reload
│   ├── ColorLoader.qml    # Theme file loader
│   ├── cava.conf          # Audio visualizer config
│   └── templates/         # Matugen templates
│       ├── hecate.json    # Shell theme template
│       ├── cava.ini       # Cava colors
│       ├── spicetify.ini  # Spotify theme
│       ├── discord.css    # Vencord/system24 theme
│       └── micro.micro    # Micro editor colorscheme
├── bar/
│   └── Bar.qml            # Main bar component
├── components/
│   ├── Workspaces.qml     # Workspace indicators
│   ├── MusicVisualizer.qml # Audio visualizer (cava)
│   ├── DateWidget.qml     # Date display
│   ├── TimeWidget.qml     # Time display
│   └── SystemWidgets.qml  # Audio + WiFi widgets
└── hecate-shell-src/      # Go CLI source
```

## Theming

HecateShell uses matugen to generate Material Design 3 color palettes. When you run `hecate wallpaper <image> -g` or `hecate theme reload`, it generates colors for:

- **HecateShell** (`theme.json`) - Shell UI colors
- **Cava** (`~/.config/cava/config`) - Audio visualizer gradient
- **Spicetify** (`~/.config/spicetify/Themes/text/color.ini`) - Spotify theme
- **Discord** (`~/.config/Vencord/themes/sys24.css`) - Vencord system24 theme
- **Micro** (`~/.config/micro/colorschemes/matugen.micro`) - Editor colorscheme

Your existing `~/.config/matugen/config.toml` templates are also applied, so any other apps you have configured will update too.

## Niri Setup

For wallpapers to appear in the overview backdrop, add these layer rules to your niri config:

```kdl
layer-rule {
    match namespace="^swww-daemon$"
    place-within-backdrop true
}

layer-rule {
    match namespace="^quickshell$"
    place-within-backdrop true
}

output "eDP-1" {
    background-color "transparent"
}
```

See [NIRI_SETUP.md](NIRI_SETUP.md) for details.

## Dependencies

- quickshell-git (Qt 6.10+, Niri support)
- cava
- pipewire + wireplumber
- matugen
- niri
- swww

## Development

```bash
# Build CLI
./build.sh

# Run in debug mode
./hecate run -rd
```

Edit any QML file and the shell will hot-reload. Edit `theme.json` and colors update within 1 second with smooth transitions.
