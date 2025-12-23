pragma Singleton
import QtQuick
import Quickshell
import Quickshell.Hyprland

// Service for Hyprland-specific functionality
QtObject {
    id: hyprlandService

    // Hyprland has no native overview mode
    property bool inOverview: false

    // Workspace model - convert Hyprland's object-based model to array-like structure
    property var workspaces: {
        var list = []
        var workspaceCount = 9 // Default to 9 workspaces

        for (var i = 1; i <= workspaceCount; i++) {
            // Find workspace in Hyprland's workspace map
            var ws = null
            if (Hyprland.workspaces) {
                for (var key in Hyprland.workspaces) {
                    if (Hyprland.workspaces[key].id === i) {
                        ws = Hyprland.workspaces[key]
                        break
                    }
                }
            }

            list.push({
                id: i,
                exists: ws !== null,
                isActive: Hyprland.focusedWorkspace ? Hyprland.focusedWorkspace.id === i : false,
                name: ws ? ws.name : ""
            })
        }
        return list
    }

    // Currently focused workspace
    property var focusedWorkspace: Hyprland.focusedWorkspace

    // Focus workspace by ID
    function focusWorkspace(id) {
        Hyprland.dispatch("workspace " + id)
        console.log("Hyprland: Switched to workspace", id)
    }

    // Focus workspace by index (0-based, converts to 1-based for Hyprland)
    function focusWorkspaceByIndex(index) {
        focusWorkspace(index + 1)
    }

    // Listen to Hyprland workspace changes to update our model
    property var workspaceConnections: Connections {
        target: Hyprland

        function onFocusedWorkspaceChanged() {
            // Trigger workspace model refresh
            hyprlandService.workspaces = hyprlandService.workspaces
        }
    }

    Component.onCompleted: {
        console.log("HyprlandService initialized")
        console.log("Hyprland focused workspace:", Hyprland.focusedWorkspace ? Hyprland.focusedWorkspace.id : "none")
    }
}
