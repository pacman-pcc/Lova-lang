#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


trap 'show_cursor' EXIT





clear_screen () {
printf "%b" "\033[2J\033[1;1H"
}

hide_cursor () {
printf "%b" "\033[?25l"
}

show_cursor () {
printf "%b" "\033[?25h"
}

matrix_text () {
local text="${1:-}"
printf "%b\n" ""
printf "%b" "\033[1;32m> "
printf "%b" "${text}"
printf "%b\n" "\033[0m"
}



hide_cursor
clear_screen

printf "%b\n" "\033[1;32m"
printf "%b\n" "01000110 01010010 01000101 01000101"
printf "%b\n" "10011001 01001111 01010110 01000001"
printf "%b\n" "\033[0m"

sleep 1

matrix_text "Wake up, Neo..."
sleep 2

matrix_text "The Matrix has you..."
sleep 2

matrix_text "Follow the white rabbit."
sleep 2

clear_screen

running=1


while [[ "$running" == 1 ]]; do
c1="ｦ ｱ ｳ ｴ ｵ ｶ ｷ ｹ ｺ ｻ ｼ ｽ ｾ ｿ"
c2="1 0 1 X Z 9 8 7 A F 0 1 1 0"
c3="ﾀ ﾂ テ ﾅ ﾆ ﾇ ﾈ ﾊ ﾋ ﾎ ﾏ ﾐ ﾑ ﾒ"
c4="0 1 0 1 0 1 1 0 0 1 1 0 0 1"
c5="ﾓ ﾔ ﾕ ﾖ ﾗ ﾘ ﾙ ﾚ ﾛ ﾜ ヰ ヱ ヲ"

printf "%b\n" "\033[0;32m${c1} \033[1;32m${c3} \033[0;32m${c5}\033[0m"
printf "%b\n" "\033[1;32m${c2} \033[0;32m${c4} \033[1;32m${c1}\033[0m"
printf "%b\n" "\033[0;32m${c5} \033[1;32m${c2} \033[0;32m${c3}\033[0m"

sleep 0.05
done
