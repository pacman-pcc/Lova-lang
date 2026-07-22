# parser.py
import re

def parse_condition(cond):
    for prefix, flag in [("is_file ", "-f"), ("is_dir ", "-d"), ("is_exist ", "-e")]:
        if cond.startswith(prefix):
            path = cond[len(prefix):].strip().strip('"\'')
            path = parse_string_inter(path)
            return f'{flag} "{path}"'

    cond = re.sub(r'\band\b', '&&', cond)
    cond = re.sub(r'\bor\b', '||', cond)
    cond = re.sub(r'\bnot\b', '!', cond)

    words = cond.split()
    processed = []
    for w in words:
        if (w.startswith('"') and w.endswith('"')) or (w.startswith("'") and w.endswith("'")):
            processed.append(w)
            continue

        w_clean = w.replace('.', '_')

        if w_clean in ("true", "false", "&&", "||", "==", "!=", ">=", "<=", ">", "<", "!") or w_clean.isdigit():
            processed.append(w)
        elif w_clean.isidentifier():
            processed.append(f'"${w_clean}"')
        elif w.startswith("$"):
            if not w.startswith('"$'):
                processed.append(f'"{w}"')
            else:
                processed.append(w)
        else:
            processed.append(w)

    return " ".join(processed)

def parse_string_inter(text):
    def replacer(match):
        var_name = match.group(1).strip().replace('.', '_')
        return f'${{{var_name}}}'
    return re.sub(r'\{([a-zA-Z0-9_\.\[\]\?]+)\}', replacer, text)

def parse_line(item):
    kind, line = item

    if kind == "EMPTY" or kind == "RAW_CMD":
        return (kind, parse_string_inter(line))

    raw_line = line.strip()

    # 1.1 Var
    if raw_line.startswith("proc "):
        content = parse_string_inter(raw_line[5:].strip())
        if "=" in content:
            var_part, val_part = content.split("=", 1)
            var_name = var_part.strip().replace('.', '_')
            val_str = val_part.strip()
            return ("ASSIGN", f"{var_name}={val_str}")
        return ("ASSIGN", content)

    # 2. Local var
    if raw_line.startswith("procloc "):
        content = parse_string_inter(raw_line[8:].strip())
        if "=" in content:
            var_part, val_part = content.split("=", 1)
            var_name = var_part.strip().replace('.', '_')
            val_str = val_part.strip()
            return ("LOCAL_ASSIGN", f"local {var_name}={val_str}")
        return ("LOCAL_ASSIGN", f"local {content}")

    # 3. Const
    if raw_line.startswith("const "):
        content = parse_string_inter(raw_line[6:].strip())
        if "=" in content:
            var_part, val_part = content.split("=", 1)
            var_name = var_part.strip().replace('.', '_')
            val_str = val_part.strip()
            return ("CONST", f"readonly {var_name}={val_str}")
        return ("CONST", f"readonly {content}")

    # 4. Del
    if raw_line.startswith("del "):
        var_name = raw_line[4:].strip().replace('.', '_')
        return ("DEL", f"unset {var_name}")

    # 2.2 Io
    if raw_line.startswith("printn ") or raw_line.startswith("print "):
        is_net = raw_line.startswith("printn ")
        prefix_len = 7 if is_net else 6
        val = raw_line[prefix_len:].strip()

        # Полное удаление всех внешних кавычек перед интерполяцией
        val = val.strip('"\'')
        val = parse_string_inter(val)

        cmd = "PRINTN" if is_net else "PRINT"
        fmt = "%b\\n" if is_net else "%b"
        return (cmd, f'printf "{fmt}" "{val}"\n')

    line_clean = parse_string_inter(raw_line)

    if line_clean.startswith("if ") and line_clean.endswith("{"):
        cond = line_clean[3:-1].strip()
        return ("IF", parse_condition(cond))

    if line_clean.startswith("elseif ") and line_clean.endswith("{"):
        cond = line_clean[7:-1].strip()
        return ("ELIF", parse_condition(cond))

    if line_clean in ("} else {", "else {"):
        return ("ELSE", "")

    if line_clean.startswith("while ") and line_clean.endswith("{"):
        cond = line_clean[6:-1].strip()
        return ("WHILE", parse_condition(cond))

    if line_clean.startswith("fdo ") and line_clean.endswith("{"):
        cond = line_clean[4:-1].strip()
        return ("UNTIL", parse_condition(cond))

    if line_clean.startswith("case ") and "do {" in line_clean:
        var = line_clean.replace("case ", "").replace("do {", "").strip()
        var = var.replace('.', '_')
        if not var.startswith("$") and not var.startswith('"'):
            var = f'"${var}"'
        return ("CASE_START", var)

    if line_clean.startswith("defer "):
        inner_cmd = line_clean[6:].strip()
        inner_kind, inner_val = parse_line(("SPECIAL", inner_cmd))
        return ("DEFER", inner_val)

    if line_clean.startswith("function ") and line_clean.endswith("{"):
        fn_header = line_clean[9:-1].strip()
        if not fn_header.endswith("()"):
            fn_header = f"{fn_header} ()"
        return ("FN_START", f"{fn_header} {{")

    if line_clean.startswith("return "):
        code = line_clean[7:].strip()
        return ("RETURN", f"return {code}")

    if line_clean == "over":
        return ("CASE_OVER", ";;")

    if line_clean == "}":
        return ("BLOCK_END", "}")

    return ("RAW", line_clean)
