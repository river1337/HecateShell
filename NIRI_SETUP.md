# Niri Configuration for HecateShell

HecateShell uses `swww` for wallpaper management. To make wallpapers appear in the backdrop (overview background) and stay stationary, add these rules to your niri config:

```kdl
// Make swww wallpaper appear in backdrop (overview background)
layer-rule {
    match namespace="^swww-daemon$"
    place-within-backdrop true
}

// Make QuickShell components visible in overview
layer-rule {
    match namespace="^quickshell$"
    place-within-backdrop true
}

// Make workspace backgrounds transparent so the backdrop shows through
output "eDP-1" {
    background-color "transparent"
}
```

## What this does:

1. **`swww-daemon` backdrop rule** - Places swww wallpapers in the backdrop
   - Wallpaper visible in overview mode
   - Stays stationary when switching workspaces
   - Shows behind all workspaces

2. **`quickshell` backdrop rule** - Makes shell components visible in overview
   - Allows the bar and other shell elements to appear properly

3. **Transparent workspace backgrounds** - Lets the backdrop wallpaper show through
   - Make `background-color` transparent in your config
   - Workspaces become transparent

## Requirements:

- `swww` daemon must be running: `swww-daemon &`
- Use `hecate wallpaper <path>` to set wallpapers
- Optionally use `-g` flag to generate theme from wallpaper colors

## Result:

✨ Consistent wallpaper across all workspaces and in overview mode!
