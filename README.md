<p align="center">
  <img src="logo.png" />
</p>

![Go Version](https://img.shields.io/badge/Go-1.24+-7A49A5?style=flat-square&logo=go&logoColor=white)
![Latest Version](https://img.shields.io/badge/version-beta-7A49A5?style=flat-square)
![License](https://img.shields.io/badge/license-GPL-5A2E85?style=flat-square)
![Platform](https://img.shields.io/badge/OS-macOS%20%2F%20Linux-7A49A5?style=flat-square&logo=linux&logoColor=white)

**Lova** - A Programming Language Built for Bash Coding Without Pain

# Features
- **Security** - from the flags `set -euo pipefail`, convenient syntax and defer security is increased many times over
- **Fast** - Transpilation of 150 lines of code takes 250-350 milliseconds
- **Simplicity** - There are no OOP, generics and extra things that need to be learned, the language is easy


# Examples

## Lova
```bash
pwd
defer ls -la
```
## Bash
```bash
#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

trap 'ls -la' EXIT

pwd
```

## Lova

```bash
proc hypers = "Lova!"
proc two_hypers = "NN!"
proc you = "Boy!"

function tur(){
    printn "\t{hypers}"
    printn "\t{two_hypers}"
    printn "\t{you}"
}

tur
```

## Bash

```bash
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
```


# Why create lova?
>Lova was created as a lighter replacement for Bash as Bash has terrible syntax and outdated design

# Documentation
-> [Guide / Documentation](guide/documentation.md)
