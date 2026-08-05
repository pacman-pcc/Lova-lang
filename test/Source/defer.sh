#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


trap 'printf "%b\n" "World!"' EXIT


printf "%b\n" "Hello"
