import QtQuick
import ".." as Shell

// Time widget component
Text {
    id: timeWidget

    property var currentTime: new Date()

    text: Qt.formatTime(currentTime, "h:mm ap")
    color: Shell.Config.textColor
    font.family: Shell.Config.fontFamily
    font.pixelSize: Shell.Config.fontSize

    // Update every second
    Timer {
        interval: 1000
        running: true
        repeat: true
        onTriggered: timeWidget.currentTime = new Date()
    }

    Component.onCompleted: {
        timeWidget.currentTime = new Date()
    }
}
