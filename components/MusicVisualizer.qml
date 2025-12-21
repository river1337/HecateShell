import QtQuick
import ".." as Shell
import Quickshell
import Quickshell.Io

// Music visualizer component with cava integration
Item {
    id: visualizer
    width: Shell.Config.visualizerWidth
    height: Shell.Config.barHeight - Shell.Config.paddingSmall * 4

    property var levels: []

    // Cava process for real-time audio visualization
    Process {
        id: cavaProcess
        running: true
        command: ["cava", "-p", Shell.Config.configPath + "/cava.conf"]

        stdout: SplitParser {
            onRead: function(data) {
                // Parse cava output: semicolon-delimited bar heights
                var bars = data.trim().split(";")
                var newLevels = []
                for (var i = 0; i < bars.length && i < Shell.Config.visualizerBarCount; i++) {
                    var value = parseInt(bars[i])
                    if (!isNaN(value)) {
                        newLevels.push(value)
                    } else {
                        newLevels.push(0)
                    }
                }
                // Ensure we have exactly visualizerBarCount bars
                while (newLevels.length < Shell.Config.visualizerBarCount) {
                    newLevels.push(0)
                }
                visualizer.levels = newLevels
            }
        }
    }

    Component.onCompleted: {
        // Initialize with zeros
        var initial = []
        for (var i = 0; i < Shell.Config.visualizerBarCount; i++) {
            initial.push(0)
        }
        visualizer.levels = initial
    }

    Row {
        anchors.bottom: parent.bottom
        anchors.horizontalCenter: parent.horizontalCenter
        spacing: Shell.Config.visualizerBarSpacing

        Repeater {
            model: Shell.Config.visualizerBarCount

            Item {
                width: (Shell.Config.visualizerWidth - (Shell.Config.visualizerBarCount - 1) * Shell.Config.visualizerBarSpacing) / Shell.Config.visualizerBarCount
                height: Shell.Config.barHeight

                Rectangle {
                    anchors.bottom: parent.bottom
                    anchors.horizontalCenter: parent.horizontalCenter
                    width: parent.width
                    height: Math.max(Shell.Config.visualizerHeight,
                                   visualizer.levels[index] !== undefined ? visualizer.levels[index] : 0)
                    color: Shell.Config.accentColor
                    radius: 1

                    Behavior on height {
                        NumberAnimation { duration: 50; easing.type: Easing.OutQuad }
                    }
                }
            }
        }
    }
}
