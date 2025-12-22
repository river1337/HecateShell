import QtQuick
import "../../.." as Shell
import QtQuick.Layouts

// Workspace indicator component
Item {
    id: workspaces

    implicitWidth: background.width
    implicitHeight: background.height

    // Background container
    Rectangle {
        id: background
        width: row.width + Shell.Config.paddingSmall * 2
        height: row.height + Shell.Config.paddingSmall * 2
        radius: height / 2
        color: Shell.Config.backgroundColorAlt

        // Sliding active indicator
        Rectangle {
            id: activeIndicator
            width: Shell.Config.workspaceItemSize
            height: Shell.Config.workspaceItemSize
            radius: width / 2
            color: Shell.Config.accentColor
            y: (background.height - height) / 2

            property int activeIndex: 0
            property int activeWorkspaceId: 1

            x: Shell.Config.paddingSmall + ((activeIndex - 1) * (Shell.Config.workspaceItemSize + Shell.Config.paddingMedium))

            Behavior on x {
                NumberAnimation { duration: 200; easing.type: Easing.OutQuad }
            }

        }

        Row {
            id: row
            anchors.centerIn: parent
            spacing: Shell.Config.paddingMedium

            Repeater {
                id: repeater
                model: niri.workspaces

                Item {
                    visible: index < 11
                    width: Shell.Config.workspaceItemSize
                    height: Shell.Config.workspaceItemSize

                    property bool isActive: model.isActive
                    property int workspaceId: model.id

                    onIsActiveChanged: {
                        if (isActive) {
                            activeIndicator.activeIndex = index
                            activeIndicator.activeWorkspaceId = workspaceId
                        }
                    }

                    Text {
                        anchors.centerIn: parent
                        text: index
                        color: model.isActive ? Shell.Config.backgroundColorAlt : Shell.Config.textColor
                        font.family: Shell.Config.fontFamily
                        font.pixelSize: Shell.Config.fontSize
                        font.bold: model.isActive

                        Behavior on color {
                            ColorAnimation { duration: 200; easing.type: Easing.OutQuad }
                        }
                    }

                    MouseArea {
                        anchors.fill: parent
                        hoverEnabled: true
                        onEntered: parent.opacity = 0.8
                        onExited: parent.opacity = 1.0
                        onClicked: {
                            niri.focusWorkspaceById(model.id)
                        }
                    }

                    Component.onCompleted: {
                        if (isActive) {
                            activeIndicator.activeIndex = index
                            activeIndicator.activeWorkspaceId = workspaceId
                        }
                    }
                }
            }
        }

        // Scroll functionality
        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: true
            onWheel: function(wheel) {
                if (wheel.angleDelta.y > 0) {
                    // Scroll up - previous workspace
                    niri.focusWorkspaceById(activeIndicator.activeWorkspaceId - 1)
                } else if (wheel.angleDelta.y < 0) {
                    // Scroll down - next workspace
                    niri.focusWorkspaceById(activeIndicator.activeWorkspaceId + 1)
                }
            }
        }
    }
}
