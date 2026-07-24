#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


trap 'rm -rf ${TMP_DIR}' EXIT

REPO="pacman-pcc/Lova-lang"
TAG="v.03"
BIN_NAME="lova"
INSTALL_DIR="/usr/local/bin"
URL="https://github.com/${REPO}/releases/download/${TAG}/${BIN_NAME}"

GREEN="\033[0;32m"
CYAN="\033[0;36m"
RED="\033[0;31m"
NC="\033[0m"

printf "%b\n" "${CYAN}==> Installing LOVA Transpiler...${NC}"

if [[ -e "/usr/bin/curl" ]]; then
printf "%b\n" "${GREEN}[✓] Found curl${NC}"
else
printf "%b\n" "${RED}[X] Error: curl is not installed.${NC}"
exit 1
fi

TMP_DIR="/tmp/lova_install"
mkdir -p ${TMP_DIR}

TARGET="/tmp/lova_install/lova"

printf "%b\n" "${CYAN}[->] Downloading binary from GitHub...${NC}"

curl -sSL ${URL} -o ${TARGET}


if [[ -f "/tmp/lova_install/lova" ]]; then
printf "%b\n" "${GREEN}[✓] Download complete.${NC}"
else
printf "%b\n" "${RED}[X] Error: Download failed.${NC}"
exit 1
fi

printf "%b\n" "${CYAN}[->] Installing to ${INSTALL_DIR}...${NC}"
chmod +x ${TARGET}

sudo mv ${TARGET} ${INSTALL_DIR}/${BIN_NAME}

if [[ -f "/usr/local/bin/lova" ]]; then
printf "%b\n" "${GREEN}=========================================${NC}"
printf "%b\n" "${GREEN}[✓] LOVA installed successfully!${NC}"
printf "%b\n" "${GREEN}[✓] Run 'lova' or 'lova run script.lova'${NC}"
printf "%b\n" "${GREEN}=========================================${NC}"
else
printf "%b\n" "${RED}[X] Installation failed.${NC}"
exit 1
fi
