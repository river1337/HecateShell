import QtQuick
import Quickshell.Io

Item {
    id: root

    signal colorsLoaded(string data)

    property string accumulatedData: ""

    function load(path) {
        accumulatedData = ""
        loader.command = ["cat", path]
        loader.running = true
    }

    Process {
        id: loader
        running: false

        stdout: SplitParser {
            onRead: function(data) {
                root.accumulatedData += data
            }
        }

        onExited: function(exitCode) {
            if (exitCode === 0 && root.accumulatedData) {
                root.colorsLoaded(root.accumulatedData)
            } else {
                console.error("Failed to read theme file, exit code:", exitCode)
            }
        }
    }
}
