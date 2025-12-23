import QtQuick
import QtQuick.Effects
import Quickshell
import Quickshell.Wayland
import "../../" as Shell

// Wallpaper module with transitions and blur support
Variants {
    id: wallpaperVariants
    model: Quickshell.screens

    delegate: PanelWindow {
        id: panel
        required property var modelData

        screen: modelData
        color: "transparent"

        // Position as background layer
        WlrLayershell.layer: WlrLayer.Background
        WlrLayershell.namespace: "hecate-wallpaper"
        WlrLayershell.exclusionMode: ExclusionMode.Ignore

        anchors {
            left: true
            right: true
            top: true
            bottom: true
        }

        Component.onCompleted: {
            console.log("Wallpaper window created - namespace:", WlrLayershell.namespace)
        }

        // Main wallpaper container
        Item {
            id: container
            anchors.fill: parent

            property string configPath: Shell.Config.wallpaperPath
            property bool useImageA: true
            property string lastPath: ""

            // Image A
            Image {
                id: imageA
                anchors.fill: parent
                fillMode: Image.PreserveAspectCrop
                asynchronous: true
                cache: false
                opacity: container.useImageA ? 1.0 : 0.0

                Behavior on opacity {
                    NumberAnimation {
                        duration: container.lastPath === "" ? 0 : Shell.Config.wallpaperDuration
                        easing.type: Easing.InOutQuad
                    }
                }
            }

            // Image B
            Image {
                id: imageB
                anchors.fill: parent
                fillMode: Image.PreserveAspectCrop
                asynchronous: true
                cache: false
                opacity: container.useImageA ? 0.0 : 1.0

                Behavior on opacity {
                    NumberAnimation {
                        duration: container.lastPath === "" ? 0 : Shell.Config.wallpaperDuration
                        easing.type: Easing.InOutQuad
                    }
                }
            }

            // Watch for wallpaper path changes
            onConfigPathChanged: {
                if (configPath === "" || configPath === lastPath) return

                console.log("Wallpaper change detected:", configPath)

                // Set the new wallpaper on the inactive image
                if (useImageA) {
                    imageB.source = configPath
                } else {
                    imageA.source = configPath
                }

                // Toggle which image is visible
                useImageA = !useImageA
                lastPath = configPath
            }

            Component.onCompleted: {
                // Initial load - set both images to the same wallpaper
                if (configPath !== "") {
                    imageA.source = configPath
                    imageB.source = configPath
                    lastPath = configPath
                    console.log("Initial wallpaper loaded:", configPath)
                }
            }

            // Blur layer for overview mode (blurs the entire container)
            // Note: Only works on Niri (Hyprland has no native overview)
            MultiEffect {
                id: blurEffect
                anchors.fill: parent
                source: container
                visible: opacity > 0
                opacity: (Shell.Config.wallpaperBlurOverview && Shell.CompositorService.inOverview) ? 1.0 : 0.0

                blur: 1.0
                blurEnabled: true
                blurMax: Shell.Config.wallpaperBlurAmount

                Behavior on opacity {
                    NumberAnimation {
                        duration: 300
                        easing.type: Easing.InOutQuad
                    }
                }
            }

        }
    }
}
