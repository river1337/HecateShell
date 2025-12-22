import QtQuick
import QtQuick.Layouts
import Quickshell
import Quickshell.Wayland
import "../.." as Shell
import "./components" as BarComponents

PanelWindow {
    id: bar

    // Bar positioning
    anchors.top: true
    anchors.left: true
    anchors.right: true

    // Bar styling
    implicitHeight: Shell.Config.barHeight
    color: Shell.Config.backgroundColor

    // LEFT SECTION: Workspaces
    BarComponents.Workspaces {
        anchors.left: parent.left
        anchors.leftMargin: Shell.Config.paddingSmall
        anchors.verticalCenter: parent.verticalCenter
    }

    // CENTER SECTION: Date, Time, and Music Visualizer
    Item {
        anchors.centerIn: parent
        implicitWidth: centerBackground.width
        implicitHeight: centerBackground.height

        Rectangle {
            id: centerBackground
            width: centerContent.width + Shell.Config.paddingLarge * 2
            height: centerContent.height + Shell.Config.paddingSmall * 2
            radius: height / 2
            color: Shell.Config.backgroundColorAlt

            Row {
                id: centerContent
                anchors.centerIn: parent
                spacing: Shell.Config.paddingLarge

                BarComponents.DateWidget {
                    anchors.verticalCenter: parent.verticalCenter
                }

                BarComponents.MusicVisualizer {
                    anchors.verticalCenter: parent.verticalCenter
                }

                BarComponents.TimeWidget {
                    anchors.verticalCenter: parent.verticalCenter
                }
            }
        }
    }

    // RIGHT SECTION: System Widgets (Audio + WiFi grouped)
    BarComponents.SystemWidgets {
        anchors.right: parent.right
        anchors.rightMargin: Shell.Config.paddingSmall
        anchors.verticalCenter: parent.verticalCenter
    }

}
