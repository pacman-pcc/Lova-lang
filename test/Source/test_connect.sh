#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
CYAN="\033[0;36m"
PURPLE="\033[0;35m"
YELLOW="\033[0;33m"
NC="\033[0m"


OS_NAME="$(uname -s)"
IS_LINUX="$([[ $(uname -s) == 'Linux' ]] && echo true || echo false)"
IS_MAC="$([[ $(uname -s) == 'Darwin' ]] && echo true || echo false)"


if [[ "${IS_LINUX}" == "true" ]]; then
printf "%b\n" "${GREEN}Running on Linux!${NC}"
fi