# About
**LOVA** is a modern, ultra-lightweight language designed to eliminate the pain of writing Bash scripts. It translate directly into safe, clean, and portable standard Bash scripts.

# Why LOVA?
- Safe (as safe as Bash gets)
  Every compiled script automatically enforces strict mode (set -euo pipefail and IFS=$'\n\t'). Catch unset variables and failing pipes instantly before they break your environment.
- ⚡ Super Lightweight
  Zero runtime overhead, no external dependencies, and a tiny footprint of just 19 keywords. It translate into plain, fast Bash that runs anywhere.
- ✨ Modern
  Clean syntax with intuitive {var} string interpolation, readable file predicates (is_file, is_dir), and Go-style defer for reliable cleanup.

# EXAMPLE
```bash
// 1. Working with Files and Safe Cleanup (defer)
function setup_environment() {
    procloc tmp_dir = $(mktemp -d)
    defer rm -rf {tmp_dir}

    if is_dir {tmp_dir} {
        printn "Created temp directory: {tmp_dir}"
    }
    return 0
}

// 2. Simple Control Flow and Interpolation
proc user = "Alex"

if is_file "./config.json" {
    printn "Welcome back, {user}! Loading config..."
} else {
    printn "Config missing! Initializing default settings for {user}."
}```
