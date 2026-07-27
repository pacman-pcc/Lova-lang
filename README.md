<p align="center">
  <img src="logo22.png" alt="Lova-2">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-7A49A5?style=flat-square&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/Version-Stable-7A49A5?style=flat-square" alt="Latest Version">
  <img src="https://img.shields.io/badge/license-GPL-5A2E85?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/OS-MacOS%20%2F%20Linux-7A49A5?style=flat-square&logo=linux&logoColor=white" alt="Platform">
  <img src="https://img.shields.io/github/v/release/pacman-pcc/Lova-lang?color=7A49A5&style=flat-square" alt="GitHub Release">
</p>

**Lova** - Transpilable programming language is designed to be as simple as possible and at the same time be ideal for scripts

# Features
- Speed - Transpilation of 100 lines to 200-250ms
- Safety - Strict Variable Safety (SVS) and set -euo pipefail ensure safety
- Synchronization - You can freely insert any bash code inside lova and it will be executed
- Simplicity - Syntax focused on the procedural paradigm, the language learns in a couple of days

# Examples

## Lova
```Bash
// ==========================================
// Along
// ==========================================

function print_neofetch() {
    // system information
    procloc user_host = $(whoami)@$(cat /etc/hostname)
    procloc os_name   = $(grep '^PRETTY_NAME=' /etc/os-release | cut -d'=' -f2 | tr -d '"')
    procloc kernel    = $(uname -r)
    procloc uptime_val = $(uptime -p)
    procloc shell_ver = $(basename "$SHELL")
    procloc cpu_model = $(grep -m 1 'model name' /proc/cpuinfo | cut -d':' -f2 | sed 's/^[ \t]*//')

    // ASCII
    printn "       /\       {user_host}"
    printn "      /  \\      ----------------------------------"
    printn "     / /\\ \\     OS: {os_name}"
    printn "    / ____ \\    Kernel: {kernel}"
    printn "   /_/    \\_\\   Uptime: {uptime_val}"
    printn "                Shell: {shell_ver}"
    printn "                CPU: {cpu_model}"
}

// Run
printn "Launching Along..."
printn ""

// Check POSIX system
if is_file "/etc/os-release" {
    print_neofetch
} else {
    printn "[ERROR] Not a standard Linux system!"
}

// defer
defer printn "Along exited cleanly."
```

## Bash
```bash
#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


trap 'printf "%b\n" "Along exited cleanly."' EXIT





print_neofetch () {

local user_host=$(whoami)@$(cat /etc/hostname)
local os_name=$(grep '^PRETTY_NAME=' /etc/os-release | cut -d'=' -f2 | tr -d '"')
local kernel=$(uname -r)
local uptime_val=$(uptime -p)
local shell_ver=$(basename "$SHELL")
local cpu_model=$(grep -m 1 'model name' /proc/cpuinfo | cut -d':' -f2 | sed 's/^[ \t]*//')


printf "%b\n" "       /\       ${user_host}"
printf "%b\n" "      /  \\      ----------------------------------"
printf "%b\n" "     / /\\ \\     OS: ${os_name}"
printf "%b\n" "    / ____ \\    Kernel: ${kernel}"
printf "%b\n" "   /_/    \\_\\   Uptime: ${uptime_val}"
printf "%b\n" "                Shell: ${shell_ver}"
printf "%b\n" "                CPU: ${cpu_model}"
}


printf "%b\n" "Launching Along..."
printf "%b\n" ""


if [[ -f "/etc/os-release" ]]; then
print_neofetch
else
printf "%b\n" "[ERROR] Not a standard Linux system!"
fi
```

## Lova
```Bash
// LOVA Script: Safe Universal Package Manager Cache Cleanup (APT, DNF, Pacman)

proc COLOR_RED = "\033[0;31m"
proc COLOR_MAGENTA = "\033[0;35m"
proc COLOR_NC = "\033[0m"

printn "{COLOR_MAGENTA}[+]{COLOR_NC} Starting universal package cache cleanup..."

// 1. APT Cache Cleanup (Debian / Ubuntu / Mint)
if is_exist "/usr/bin/apt-get" {
    printn "{COLOR_MAGENTA}[*]{COLOR_NC} Cleaning APT cache..."
    // Clean obsolete packages and partial downloads safely
    apt-get clean
    apt-get autoclean
}

// 2. DNF Cache Cleanup (Fedora / RHEL / Rocky Linux)
if is_exist "/usr/bin/dnf" {
    printn "{COLOR_MAGENTA}[*]{COLOR_NC} Cleaning DNF cache..."
    // Clean expired and old cached metadata/packages
    dnf clean expire-cache
    dnf clean packages
}

// 3. Pacman Cache Cleanup (Arch Linux / Manjaro / EndeavourOS)
if is_exist "/usr/bin/paccache" {
    printn "{COLOR_MAGENTA}[*]{COLOR_NC} Cleaning Pacman cache via paccache..."
    // Keep the 2 most recent versions and clear completely uninstalled packages
    paccache -rk 2
    paccache -ruk 0
} else if is_exist "/usr/bin/pacman" {
    printn "{COLOR_RED}[*] Cleaning Pacman cache via standard pacman...{COLOR_NC}"
    // Safely remove cached versions of uninstalled packages
    pacman -Sc --noconfirm
}

printn "{COLOR_MAGENTA}[✓]{COLOR_NC} All available package caches cleaned successfully!"
```

## Bash
```Bash
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
```

# Speed in 10 line
```
$ lova -r new.lova
Translate: new.lova :: new.sh..
LOVA: Time compile: 1.440137ms
	Lova!
	NN!
	Boy!
```

# Who Needs the Lova Programming Language
>DevOps, System Administrators, Clouds (due to the weight of the transpiler of only 3-4MB), and those who are annoyed by Bash

# Documentation
-> [Guide / Documentation](guide/documentation.md)

# Installing
```bash
bash <(curl -sSL https://raw.githubusercontent.com/pacman-pcc/Lova-lang/main/install.sh)
```
