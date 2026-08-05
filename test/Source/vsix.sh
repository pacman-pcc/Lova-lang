#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

REPO="pacman-pcc/Lova-Snippets"
BIN_NAME="main/lova-1.0.0.vsix"
URL="https://raw.githubusercontent.com/${REPO}/${BIN_NAME}"

GREEN="\033[0;32m"
CYAN="\033[0;36m"
RED="\033[0;31m"
NC="\033[0m"

printf "%b\n" "${CYAN}==>${NC} Installing LOVA VSIX..."

if [[ -e "/usr/bin/curl" ]]; then
printf "%b\n" "${GREEN}[+]${NC} Found curl"
else
printf "%b\n" "${RED}[X]${NC} Error: curl is not installed."
exit 1
fi

printf "%b\n" "${CYAN}==>${NC} Downloading VSIX from GitHub..."

curl -sSL ${URL} -o $HOME/lova-1.0.0.vsix

if [[ -f "$HOME/lova-1.0.0.vsix" ]]; then
printf "%b\n" "${GREEN}[+]${NC} Download complete."
else
printf "%b\n" "${RED}[X]${NC} Error: Download failed."
exit 1
fi

printf "%b\n" "${CYAN}==>${NC} Installing VSIX to Visual Studio Code..."

if [[ -e "/usr/bin/code" ]]; then
exec code --install-extension $HOME/lova-1.0.0.vsix
printf "%b\n" "${GREEN}\t==>${NC} LOVA VSIX installed successfully!"
else
printf "%b\n" "${RED}[X]${NC} Error: VS Code CLI ('code') is not found."
exit 1
fi