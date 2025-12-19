import QtQuick
import ".." as Shell

// Date widget component
Text {
    id: dateWidget

    property var currentDate: new Date()

    text: Qt.formatDate(currentDate, "MMM. d")
    color: Shell.Config.textColor
    font.family: Shell.Config.fontFamily
    font.pixelSize: Shell.Config.fontSize

    // Update every minute
    Timer {
        interval: 60000 // 1 minute
        running: true
        repeat: true
        onTriggered: dateWidget.currentDate = new Date()
    }

    Component.onCompleted: {
        dateWidget.currentDate = new Date()
    }
}
