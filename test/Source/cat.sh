#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

frame=0

while [[ true ]]; do
printf "%b" "\033[H\033[J"

case "$frame" in
"0")
printf "%b\n" " /\\_/\\"
printf "%b\n" "( o.o )"
printf "%b\n" " > ^ < "
;;
*)
printf "%b\n" " /\\_/\\"
printf "%b\n" "( -.- )"
printf "%b\n" " > ^ < "
;;
esac

sleep 1
if [[ "$frame" == 0 ]]; then
frame=1
else
frame=0
fi
done
