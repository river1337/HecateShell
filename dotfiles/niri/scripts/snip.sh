#!/usr/bin/env bash

# temp files
TMP_FULL=$(mktemp /tmp/snip_full_XXXX.png)
TMP_CROP=$(mktemp /tmp/snip_crop_XXXX.png)

# take fullscreen screenshot
grim "$TMP_FULL"

# show overlay (fullscreen, fully opaque)
imv -f "$TMP_FULL" &
OVERLAY_PID=$!

# tiny sleep to ensure overlay is visible
sleep 0.1

# select area
GEOM=$(slurp -d -c "#ffffff" -b "#00000099")

# kill overlay
kill "$OVERLAY_PID"

# if nothing selected, exit
[[ -z "$GEOM" ]] && { rm "$TMP_FULL" "$TMP_CROP"; exit; }

# reformat geometry for convert: slurp -> "WxH+X+Y"
read X Y W H <<< $(echo $GEOM | sed 's/[ ,x]/ /g')
CROP="${W}x${H}+${X}+${Y}"

# crop
convert "$TMP_FULL" -crop "$CROP" +repage "$TMP_CROP"

# copy to clipboard
wl-copy < "$TMP_CROP"

# cleanup
rm "$TMP_FULL" "$TMP_CROP"
