# Custom QuickShell for Niri/Wayland

A modular, customizable shell built with QuickShell for Arch Linux running Niri (Wayland).

## Project Structure

```
.
├── shell.qml              # Main entry point
├── qmldir                 # QML module definitions
├── config/
│   ├── Config.qml         # Central config (hot-reloadable)
│   └── theme.conf         # Theme configuration file
├── bar/
│   └── Bar.qml            # Main bar component
├── components/
│   ├── Workspaces.qml     # Workspace indicator
│   ├── MusicVisualizer.qml # Audio visualizer
│   ├── DateWidget.qml     # Date display
│   ├── TimeWidget.qml     # Time display
│   ├── AudioWidget.qml    # Volume control widget
│   └── WifiWidget.qml     # Network status widget
└── modules/               # Future modules
    ├── wallpapers/
    ├── theme/
    ├── appfinder/
    ├── settings/
    └── notifications/
```

## Features

### Bar
- **Left**: Workspace indicators (5 workspaces)
- **Center**: Date | Music Visualizer | Time
- **Right**: Audio widget | WiFi widget

### Hot Reloading
Edit `config/Config.conf` or any QML file and the shell will automatically reload!

### Customization
All styling is centralized in `config/Config.qml`:
- Colors (Matugen color generation)
- Fonts (JetBrains Mono by default)
- Spacing, padding, sizes
- Widget-specific settings

## Current Status

### Working
- ✓ Bar layout with all widgets
- ✓ Hot-reloading configuration
- ✓ Modular component structure
- ✓ Date/Time widgets (live updating)
- ✓ Music visualizer (placeholder animation)
- ✓ Audio widget (placeholder)
- ✓ WiFi widget (placeholder)
- ✓ Workspace indicators (placeholder)

### TODO
- [ ] Add wallpaper manager
- [ ] Add theme selector with matugen
- [ ] Add app finder/launcher
- [ ] Add settings menu
- [ ] Add notifications system

## Development

### Testing
Click on widgets to test interactions:
- Workspaces: Click to "switch" (console log)
- Music Visualizer: Click to toggle animation
- Audio: Click to toggle mute
- WiFi: Click to toggle connection state

### Adding New Modules
1. Create new directory in `modules/`
2. Create QML component
3. Add to `qmldir`
4. Import in `shell.qml`

### Customizing Theme
Edit `config/Config.qml` to change:
- Colors
- Fonts
- Sizes
- Icons

## Notes

- Using Unicode icons for now (Nerd Fonts recommended)
- Install `JetBrains Mono Nerd Font` or change font in Config.qml
- Bar height is 35px by default (adjustable in Config.qml)
