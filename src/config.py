# config.py
#
# safety config Bash

BASH_HEADERS = """#!/bin/bash
set -euo pipefail
IFS=$'\\n\\t'
"""

# operations in bash
OPERATORS = {
    "==": "-eq",
    "!=": "-ne",
    ">": "-gt",
    "<": "-lt",
    ">=": "-ge",
    "<=": "-le",
    "and": "&&",
    "or": "||",
}

# keywords in OREZ
#

KEYWORDS = {
    "proc",
    "procloc",
    "const",
    "del",
    "printn",
    "print",
    "if",
    "elseif",
    "else",
    "while",
    "fdo",
    "case",
    "over",
    "function",
    "return",
    "defer",
}
