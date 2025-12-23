import Quickshell
import QtQuick
import "." as Shell
import Niri 0.1
import Quickshell.Hyprland

ShellRoot {
    id: shell

    // Niri compositor (only used if on Niri)
    Niri {
        id: niri
        Component.onCompleted: {
            if (Shell.CompositorService.isNiri) {
                connect()
            }
        }

        onConnected: console.info("Connected to Niri")
        onErrorOccurred: function(error) {
            console.error("Niri error:", error)
        }
    }

    // Load modules
    Wallpaper {}
    Bar {}

    // Future modules can be added here:
    // Notifications {}
    // AppFinder {}
    // Settings {}

    Component.onCompleted: {
        console.log("Shell initialized!")
        console.log("Hot-reloading enabled - edit theme.json to update theme")

        // Wire up compositor-specific services
        if (Shell.CompositorService.isNiri) {
            Shell.NiriService.niriObject = niri
            Shell.NiriService.workspaces = Qt.binding(() => niri.workspaces)
            Shell.CompositorService.niriService = Shell.NiriService
        } else if (Shell.CompositorService.isHyprland) {
            Shell.CompositorService.hyprlandService = Shell.HyprlandService
        }

        // Force Config singleton to initialize
        var _ = Shell.Config.barHeight
        console.log("Config initialized")
    }
}
