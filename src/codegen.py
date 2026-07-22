# codegen.py
from config import BASH_HEADERS

def gen_bash(parsed_items):
    out = [BASH_HEADERS]
    block_stack = []
    defer_stack = []

    for kind, val in parsed_items:
        if kind in ("RAW_CMD", "EMPTY", "ASSIGN", "LOCAL_ASSIGN", "CONST", "DEL", "PRINTN", "PRINT", "RAW"):
            out.append(val)

        elif kind == "IF":
            cond = " ".join([f'"${w}"' if w.isidentifier() and w not in ("true", "false") else w for w in val.strip().replace(" and ", " && ").replace(" or ", " || ").split()])
            out.append(f"if [[ {cond} ]]; then")
            block_stack.append("if")

        elif kind == "ELIF":
            cond = " ".join([f'"${w}"' if w.isidentifier() and w not in ("true", "false") else w for w in val.strip().replace(" and ", " && ").replace(" or ", " || ").split()])
            out.append(f"elif [[ {cond} ]]; then")

        elif kind == "ELSE":
            out.append("else")

        elif kind == "WHILE":
            out.append(f"while [[ {val} ]]; do")
            block_stack.append("loop")

        elif kind == "UNTIL":
            out.append(f"until [[ {val} ]]; then")
            block_stack.append("loop")

        elif kind == "DEFER":
            defer_stack.insert(0, val)
            continue

        elif kind == "CASE_START":
            out.append(f"case \"${val}\" in")
            block_stack.append("case")

        elif kind == "CASE_OVER":
            out.append(";;")

        elif kind == "FN_START":
            out.append(val)

        elif kind == "RETURN":
            out.append(f"   {val}")

        elif kind == "BLOCK_END":
            if not block_stack:
                out.append("}")  # Обычная закрывающая скобка (например, для функций)
                continue

            last_block = block_stack.pop()
            if last_block == "if":
                out.append("fi")
            elif last_block == "loop":
                out.append("done")
            elif last_block == "case":
                out.append("esac")

    if defer_stack:
        defer_string = "; ".join(defer_stack)
        out.insert(1, f"\ntrap '{defer_string}' EXIT\n")

    return "\n".join(out)
