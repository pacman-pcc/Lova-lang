#!/bin/bash
set -euo pipefail
IFS=$'\n\t'



COLOR_RED="\033[0;31m"
COLOR_MAGENTA="\033[0;35m"
COLOR_NC="\033[0m"

printf "%b\n" "${COLOR_MAGENTA}[+]${COLOR_NC} Starting universal package cache cleanup..."


if [[ -e "/usr/bin/apt-get" ]]; then
printf "%b\n" "${COLOR_MAGENTA}[*]${COLOR_NC} Cleaning APT cache..."

apt-get clean
apt-get autoclean
fi


if [[ -e "/usr/bin/dnf" ]]; then
printf "%b\n" "${COLOR_MAGENTA}[*]${COLOR_NC} Cleaning DNF cache..."

dnf clean expire-cache
dnf clean packages
fi


if [[ -e "/usr/bin/paccache" ]]; then
printf "%b\n" "${COLOR_MAGENTA}[*]${COLOR_NC} Cleaning Pacman cache via paccache..."

paccache -rk 2
paccache -ruk 0
elif [[ -e "/usr/bin/pacman" ]]; then
printf "%b\n" "${COLOR_RED}[*] Cleaning Pacman cache via standard pacman...${COLOR_NC}"

pacman -Sc --noconfirm
fi

printf "%b\n" "${COLOR_MAGENTA}[✓]${COLOR_NC} All available package caches cleaned successfully!"
