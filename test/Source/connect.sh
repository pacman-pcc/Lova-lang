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


printf "%b\n" "${RED}Hello,${NC}${GREEN} World!${NC}"
printf "%b\n" "${PURPLE}This is a test of the colors library.${NC}"
printf "%b\n" "${BLUE}Good${NC}${YELLOW} day!${NC}"