#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

hypers="Lova!"
two_hypers="NN!"
you="Boy!"

tur () {
printf "%b\n" "\t${hypers}"
printf "%b\n" "\t${two_hypers}"
printf "%b\n" "\t${you}"
}

tur
