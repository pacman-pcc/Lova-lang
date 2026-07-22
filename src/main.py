# main.py
# main enter

import sys
import os
import glob
from lexer import lex
from parser import parse_line
from codegen import gen_bash
import subprocess

def translate_file(oz_path):
    if not oz_path.endswith(".lova"):
        print(f"LOVA: File {oz_path} not found extension .lova")
        return

    sh_path = oz_path[:-5] + ".sh"

    print(f"Translate: {oz_path} :: {sh_path}..")
    try:
        with open(oz_path, "r", encoding="utf-8") as f:
            code = f.read()

        raw_lines = lex(code)

        parsed_items = [parse_line(item) for item in raw_lines]

        bash_code = gen_bash(parsed_items)

        with open(sh_path, "w", encoding="utf-8") as f:
            f.write(bash_code)

        os.chmod(sh_path, 0o775)
        print("LOVA: Ready ::")
    except Exception as e:
        print(f"LOVA: Error build :: {e}")
        return None
def main():
    if len(sys.argv) > 1:
        target = sys.argv[1]

        if target == "run":
            target = sys.argv[2] if len(sys.argv) > 2 else "main.lova"
            if os.path.isfile(target):
                translate_file(target)
                sh_file = target[:-5] + ".sh"
                subprocess.run(["bash", sh_file])
            else:
                print(f"LOVA: {target} not found.")

        elif os.path.isfile(target):
            translate_file(target)
        else:
            print(f"LOVA: {target} not found.")
    else:
        oz_files = glob.glob("*.lova")
        if not oz_files:
            print("LOVA: *.lova files not found in directory")

        for oz_file in oz_files:
            translate_file(oz_file)

if __name__ == "__main__":
    main()
