# LOVA Language Specification & Reference Manual

Welcome to the official language reference for **LOVA**—a modern, ultra-lightweight scripting language and **transpiler** designed to target POSIX Bash.

LOVA bridges the gap between raw shell speed and modern programming ergonomic standards. It introduces block scoping, memory cleanup mechanics, strict runtime safety, and intuitive file predicates, transpiling directly into clean, standalone Bash scripts.


## 🎯 Language Philosophy

Writing complex shell scripts in standard Bash often leads to bug-prone, hard-to-read code. Unset variables silently pass, pipe failures go unnoticed, and file testing requires cryptic flags (`-f`, `-d`, `-e`).

LOVA solves these issues by adhering to four core principles:
* **Zero-Config Safety**: All transpiled scripts run under strict execution safety (`set -euo pipefail` and `IFS=$'\n\t'`).
* **Minimal Core**: A tiny keyword footprint (19 keywords total) ensuring extreme simplicity and low cognitive load.
* **Modern Syntax**: Clean string interpolation, explicit scope definitions (`proc` vs `procloc`), and block-level syntax (`{ ... }`).
* **Deterministic Cleanup**: Native Go-style `defer` statement transpiled to reliable POSIX traps.

---

## ⚡ Transpilation & Execution Model

The LOVA toolchain operates as a source-to-source transpiler (`lova` CLI binary). It reads `.lova` source files, parses tokens, applies safety wrappers, and generates standard `.sh` output.

```bash
[ Source Code: app.lova ]  --->  ( LOVA Transpiler / lova )  --->  [ Output: app.sh ]
                                                                        |
                                                                ( Executed via Bash )
```

### CLI Usage

```bash
# Transpile all .lova files in the current workspace to .sh
lova

# Transpile and execute a specific .lova file immediately
lova run script.lova
```

---

## 🔑 Keywords & Grammar Reference

LOVA reserves **19 keywords**. Variables, function names, and custom identifiers must not collide with these names:

| Keyword | Category | Description |
| :--- | :--- | :--- |
| `function` | Declaration | Defines a procedure block |
| `return` | Control Flow | Returns an exit code from a function |
| `defer` | Cleanup | Registers a deferred action upon script exit |
| `proc` | Variable | Declares a global or script-scoped variable |
| `procloc` | Variable | Declares a local variable inside a function scope |
| `const` | Variable | Declares a read-only variable |
| `del` | Variable | Unsets/deletes a variable from memory |
| `printn` | I/O | Prints text to standard output with a newline |
| `print` | I/O | Prints text to standard output without a newline |
| `if` | Control Flow | Conditional evaluation branch |
| `elseif` | Control Flow | Secondary conditional evaluation branch |
| `else` | Control Flow | Fallback evaluation branch |
| `while` | Control Flow | Loop condition evaluation |
| `fdo` | Control Flow | Loop execution block indicator |
| `case` | Control Flow | Pattern matching construct |
| `over` | Control Flow | Pattern matching block terminator |
| `is_file` | Predicate | Evaluates if path exists and is a regular file |
| `is_dir` | Predicate | Evaluates if path exists and is a directory |
| `is_exist` | Predicate | Evaluates if path or file exists |

---

## Variables

```bash
// proc [name] = [count]
proc age = 22
```

## Constant

```bash
const name = "Tom"
```

## if/elseif/else

```bash
if age > 22 {
    printn "Big"
} else if {
    printn "22"
} else {
    printn "Small"
}
```

## defer

```bash
pwd
defer ls -la
```


## Loop

```bash
proc a = 0
while a > 2 {
    ls -la
    a = a + 1
}```

```bash
proc is_happy = false
fdo is_happy {
    print "Lalala"
}```

## Delete variable

```bash
proc a = 2
del a
```

## Function

```bash
function calc(){
    proc a = 2
    proc b = 2

    printn "{a} + {b}"
    // return command )))
}
```

## case
```bash
proc action = "Start"

case action do {
    "start":
        printn "start.."
        over
    "stop" | "pause":
        printn "stoping..."
        over
    _:
        printn "Error!"
        over
}
```

## simplecase

`
is_dir -> -d in bash
is_file -> -f in bash
is_exist
`
