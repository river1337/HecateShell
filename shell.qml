import Quickshell
import QtQuick
import "." as Shell
import Niri 0.1

ShellRoot {
    id: shell

    // Niri
    Niri {
        id: niri
        Component.onCompleted: connect()

        onConnected: console.info("Connected to niri")
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

        // Force Config singleton to initialize
        var _ = Shell.Config.barHeight
        console.log("Config initialized")
    }
}
