pragma Singleton
import QtQuick
import Quickshell

// Service to detect compositor and provide unified interface
QtObject {
    id: compositorService

    // Compositor detection via environment variables
    property bool isHyprland: Quickshell.env("HYPRLAND_INSTANCE_SIGNATURE") !== ""
    property bool isNiri: Quickshell.env("NIRI_SOCKET") !== ""
    property string compositor: isHyprland ? "hyprland" : isNiri ? "niri" : "unknown"

    // Unified properties (delegated to active service)
    property bool inOverview: isNiri ? niriService.inOverview : false
    property var workspaces: isNiri ? niriService.workspaces : hyprlandService.workspaces
    property var focusedWorkspace: isNiri ? niriService.focusedWorkspace : hyprlandService.focusedWorkspace

    // Reference to compositor-specific services
    property var niriService: null
    property var hyprlandService: null

    // Unified methods
    function focusWorkspace(id) {
        if (isNiri && niriService) {
            niriService.focusWorkspace(id)
        } else if (isHyprland && hyprlandService) {
            hyprlandService.focusWorkspace(id)
        }
    }

    function focusWorkspaceByIndex(index) {
        if (isNiri && niriService) {
            niriService.focusWorkspaceByIndex(index)
        } else if (isHyprland && hyprlandService) {
            hyprlandService.focusWorkspaceByIndex(index)
        }
    }

    Component.onCompleted: {
        console.log("CompositorService initialized")
        console.log("Detected compositor:", compositor)

        if (compositor === "unknown") {
            console.warn("WARNING: Unknown compositor! Shell may not function correctly.")
            console.warn("Set NIRI_SOCKET or HYPRLAND_INSTANCE_SIGNATURE environment variable.")
        }
    }
}
