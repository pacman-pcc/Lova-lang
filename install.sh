#!/bin/bash
set -euo pipefail
IFS=$'\n\t'



repo="https://github.com/pacman-pcc/Lova-lang"
install_dir="${HOME}/.local/bin"
binary_name="lova"
download_url="0"
os_type="$(uname -s | tr -d '\r\n')"

green="\033[0;32m"
red="\033[0;31m"
reset="\033[0m"

cd ~

printf "%b\n" "${green}Installing LOVA..${reset}"



mkdir -p ${install_dir}

printf "%b\n" "${green}Download binary..${reset}"

curl -fsSL ${download_url} -o "${install_dir}/${binary_name}"

printf "%b\n" "${os_type}"


if [[ "$os_type" != "Linux" && "$os_type" != "Darwin" ]]; then
printf "%b\n" "${red}Your OS is not suitable for the language${reset}"

else
printf "%b\n" "${green}Good${reset}"

fi

chmod +x "${install_dir}/${binary_name}"


printf "%b\n" "Lova installed in ${install_dir}/${binary_name}"


printf "%b\n" ""
