#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


trap 'ls -la' EXIT

pwd