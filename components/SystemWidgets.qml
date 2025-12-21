import QtQuick
import ".." as Shell
import Quickshell
import Quickshell.Io

// System widgets container - ONE background for both audio and wifi
Item {
    id: systemWidgets

    implicitWidth: background.width
    implicitHeight: background.height

    // PipeWire/WirePlumber for audio
    Process {
        id: volumeGetter
        running: true
        command: ["wpctl", "get-volume", "@DEFAULT_AUDIO_SINK@"]

        property int volume: 0
        property bool isMuted: false
        property int previousVolume: 50  // Store previous non-zero volume

        stdout: SplitParser {
            onRead: function(data) {
                var output = data.trim()
                // Format: "Volume: 0.50" or "Volume: 0.50 [MUTED]"
                var parts = output.split(" ")
                if (parts.length >= 2) {
                    var vol = parseFloat(parts[1]) * 100
                    volumeGetter.volume = Math.round(vol)
                    volumeGetter.isMuted = output.includes("[MUTED]")

                    // Store previous non-zero volume for toggle functionality
                    if (vol > 0) {
                        volumeGetter.previousVolume = Math.round(vol)
                    }
                }
            }
        }
    }

    // Poll volume every second
    Timer {
        interval: 1000
        running: true
        repeat: true
        onTriggered: volumeGetter.running = true
    }

    // Check internet connectivity
    Process {
        id: networkChecker
        running: true
        command: ["ping", "-c", "1", "-W", "1", "1.1.1.1"]

        property bool isOnline: false

        onExited: function(exitCode) {
            networkChecker.isOnline = (exitCode === 0)
        }
    }

    // Poll network every 5 seconds
    Timer {
        interval: 5000
        running: true
        repeat: true
        onTriggered: networkChecker.running = true
    }

    Rectangle {
        id: background
        width: content.width + Shell.Config.paddingLarge * 2
        height: Shell.Config.workspaceItemSize + Shell.Config.paddingSmall * 2
        radius: height / 2
        color: Shell.Config.backgroundColorAlt

        Row {
            id: content
            anchors.centerIn: parent
            spacing: Shell.Config.paddingLarge

            // Audio Widget
            Row {
                spacing: Shell.Config.paddingSmall

                Text {
                    text: (volumeGetter.isMuted || volumeGetter.volume === 0) ? Shell.Config.iconVolumeMuted : Shell.Config.iconVolume
                    color: (volumeGetter.isMuted || volumeGetter.volume === 0) ? Shell.Config.textColorDim : Shell.Config.accentColor
                    font.family: Shell.Config.fontFamily
                    font.pixelSize: Shell.Config.fontSizeLarge
                    anchors.verticalCenter: parent.verticalCenter

                    Behavior on color {
                        ColorAnimation { duration: 200; easing.type: Easing.OutQuad }
                    }
                }

                Text {
                    text: volumeGetter.volume + "%"
                    color: (volumeGetter.isMuted || volumeGetter.volume === 0) ? Shell.Config.textColorDim : Shell.Config.textColor
                    font.family: Shell.Config.fontFamily
                    font.pixelSize: Shell.Config.fontSize
                    anchors.verticalCenter: parent.verticalCenter

                    Behavior on color {
                        ColorAnimation { duration: 200; easing.type: Easing.OutQuad }
                    }
                }
            }

            // WiFi Widget
            Row {
                spacing: Shell.Config.paddingSmall

                Text {
                    text: networkChecker.isOnline ? Shell.Config.iconWifiOn : Shell.Config.iconWifiOff
                    color: networkChecker.isOnline ? Shell.Config.accentColor : Shell.Config.textColorDim
                    font.family: Shell.Config.fontFamily
                    font.pixelSize: Shell.Config.fontSizeLarge
                    anchors.verticalCenter: parent.verticalCenter

                    Behavior on color {
                        ColorAnimation { duration: 200; easing.type: Easing.OutQuad }
                    }
                }

                Text {
                    text: networkChecker.isOnline ? "Online" : "Offline"
                    color: networkChecker.isOnline ? Shell.Config.textColor : Shell.Config.textColorDim
                    font.family: Shell.Config.fontFamily
                    font.pixelSize: Shell.Config.fontSize
                    anchors.verticalCenter: parent.verticalCenter

                    Behavior on color {
                        ColorAnimation { duration: 200; easing.type: Easing.OutQuad }
                    }
                }
            }
        }

        // Audio controls
        MouseArea {
            x: Shell.Config.paddingLarge
            y: 0
            width: content.children[0].width
            height: parent.height
            hoverEnabled: true
            onClicked: {
                // Toggle between 0% and previous volume
                var targetVolume = (volumeGetter.volume === 0) ? volumeGetter.previousVolume : 0
                var volProcess = Qt.createQmlObject('
                    import Quickshell.Io
                    Process {
                        running: true
                        command: ["wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", "' + (targetVolume / 100.0) + '"]
                    }
                ', systemWidgets)
                volProcess.running = true
                Qt.callLater(function() { volumeGetter.running = true })
            }
            onWheel: function(wheel) {
                var delta = wheel.angleDelta.y > 0 ? "5%+" : "5%-"
                var volProcess = Qt.createQmlObject('
                    import Quickshell.Io
                    Process {
                        running: true
                        command: ["wpctl", "set-volume", "@DEFAULT_AUDIO_SINK@", "' + delta + '", "-l", "1.0"]
                    }
                ', systemWidgets)
                volProcess.running = true
                Qt.callLater(function() { volumeGetter.running = true })
            }
        }

        // WiFi area (no click action)
        // MouseArea removed as requested
    }
}
