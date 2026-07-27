#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

hello () {
local name="Lova!"
printf "%b\n" "Hello, ${name}"
}



hello
