# lexer.py
# Is this lexer OZER

from config import KEYWORDS

def tokenize(line):
    line = line.strip()

    if not line or line.startswith("//"):
        return [("EMPTY", "")]

    if " //" in line:
        line = line.split(" //")[0].strip()

    words = line.split()
    first_word = words[0]

    if first_word in KEYWORDS or first_word in ["}", "fdo"]:
        return [("SPECIAL", line)]

    return [("RAW_CMD", line)]

def lex(code):
    lines = code.splitlines()
    parsed_lines = []

    for line in lines:
        parsed_lines.extend(tokenize(line))

    return parsed_lines
