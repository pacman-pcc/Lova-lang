#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

angle=0

while [[ true ]]; do
printf "%b" "\033[H\033[J"

case "$angle" in
"0")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /    .    \\"
printf "%b\n" "    |    (O)    |"
printf "%b\n" "     \\    '    /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"45")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /   .     \\"
printf "%b\n" "    |   (O)     |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"90")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /         \\"
printf "%b\n" "    |  (O)      |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"135")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /         \\"
printf "%b\n" "    | (O)       |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"180")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /         \\"
printf "%b\n" "    |     ~     |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"225")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /         \\"
printf "%b\n" "    |       (O) |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
"270")
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /         \\"
printf "%b\n" "    |      (O)  |"
printf "%b\n" "     \\         /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
*)
printf "%b\n" "        _..._"
printf "%b\n" "      .'     '."
printf "%b\n" "     /     .   \\"
printf "%b\n" "    |     (O)   |"
printf "%b\n" "     \\    '    /"
printf "%b\n" "      '.     ."
printf "%b\n" "        '---"
;;
esac

sleep 0.3
angle=$((angle + 45))
if [[ "$angle" == 360 ]]; then
angle=0
fi
done
