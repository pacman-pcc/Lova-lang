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


