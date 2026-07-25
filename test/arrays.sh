#!/bin/bash
set -euo pipefail
IFS=$'\n\t'


my_array=("Apple", "Banana", "Orange")


printf "%b\n" "All elements: ${my_array[@]}"


my_array+=("Pear", "Kiwi")
printf "%b\n" "After append: ${my_array[@]}"


unset 'my_array[1]'
printf "%b\n" "After deleting index 1: ${my_array[@]}"
