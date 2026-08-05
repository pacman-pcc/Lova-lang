#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

Name="Tom"

if [[ "$Name" == "NN" ]]; then
printf "%b\n" "Hello, sir!"
elif [[ "$Name" == "Lola" ]]; then
printf "%b\n" "Hello, sis!"
else
printf "%b\n" "Permission denied!"
fi
