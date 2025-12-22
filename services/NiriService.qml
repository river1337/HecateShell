pragma Singleton
import QtQuick
import Quickshell
import Quickshell.Io

// Service to track Niri compositor state
QtObject {
    id: niriService

    // Track if we're in overview mode
    property bool inOverview: false

    // Niri IPC socket path
    property string socketPath: Quickshell.env("NIRI_SOCKET")

    // Process to listen to Niri events
    property var eventProcess: Process {
        id: eventListener
        running: niriService.socketPath !== ""
        command: ["niri", "msg", "--json", "event-stream"]

        stdout: SplitParser {
            onRead: function(data) {
                try {
                    var event = JSON.parse(data.trim())
                    handleNiriEvent(event)
                } catch (e) {
                    console.warn("Failed to parse Niri event:", e)
                }
            }
        }

        onExited: function(exitCode, exitStatus) {
            if (exitCode !== 0) {
                console.error("Niri event stream exited with code:", exitCode)
            }
        }
    }

    function handleNiriEvent(event) {
        // Handle overview opened/closed event
        if (event.OverviewOpenedOrClosed !== undefined) {
            var wasInOverview = inOverview
            inOverview = event.OverviewOpenedOrClosed.is_open

            if (wasInOverview !== inOverview) {
                console.log("Niri overview:", inOverview ? "opened" : "closed")
            }
        }
    }

    Component.onCompleted: {
        if (socketPath === "") {
            console.warn("NIRI_SOCKET not set - overview detection disabled")
        } else {
            console.log("NiriService initialized - listening for events")
        }
    }
}
