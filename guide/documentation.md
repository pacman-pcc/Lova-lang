# Lova 
Lova - Why is it needed? What kind of language? How does it work? I will explain

## 1. Transpile work
What is transpilation? Transpilation is when you take your code and translate it into another, is this language considered just a shell? no.. many languages work this way, for example, TypeScript, Nim, Vlang and even the initial C++

Here's how it works

```Lova-code --- Lovac (lovatranspiler) --- *.sh```

Are there any disadvantages in transpilation, maybe the interpretation is better? Not at all, transpiling to Bash gives you the ability to insert, also gives you full sync since Bash is almost everywhere

## 2. Ready

Well, now that we have understood the principle of the transpiler, we can start learning, the first one is downloaded by the command in the README 
lova then type `lova -v` to check the version'

`lova` => Compiles all *.lova into *.sh that is in the folder

`lova -v` => shows the version

`lova -r <file.lova` => transpils the specific file and runs it immediately 

`lova <file.lova>` => just compiles

## 3. Syntax

### Procedure (variables)
Well, variables or procedures (call it what you want) The keyword in the variables is 'proc', followed by the name and through = assignment, you do not need to specify types since the language is dynamic

```bash
// example
proc age = 22
age = age + 2 // 24
age = age - 2 // 20
age = age * 2 // 44
age = age / 2 // 11
```

### I/O
Here we have the I/O

`printn "text"` => print and newline (\n)

`print "text"` => only print

`{}` in `printn/print` put a variable inside

```bash
proc age = 19

printn "Age: {age}" // Age: 19
```

To accept command-line arguments, there is $1 for the first argument, $2 for the second argument

```bash
proc input = $1
```

BUT what the... unbound variables? and SVS? This is protection built into security, but in order not to leave it alone, you can use, for example, this

```bash
proc inputable = ${1:-} // ${1:-default}
```

### IF/ELSEIF/ELSE
The most interesting thing is that we do not hesitate 

if - if the condition is equal to.. 

elseif - if the condition is equal to another 

else - if the condition is neither if nor elseif

```bash
proc Name = "Tom"

if Name == "NN" {
    printn "Hello, sir!"
} else if Name == "Lola" {
    printn "Hello, sis!"
} else {
    printn "Permission denied!"
}
```


### Function

Functions are a convenient way to use the same thing without writing a bunch of code

```bash
function hello(){
    printn "Hello, Lova!"
}

hello // if args to => hello arg arg2
```

and procloc

procloc = this is a local variable inside the function

```bash
function hello(){
    procloc name = "Lova!"
    printn "Hello, {name}"
}

hello // if args to => hello arg arg2
```

```
$ lova -r func.lova
LOVA: Translate: func.lova :: func.sh..
LOVA: Time compile: 1.096391ms
Hello, Lova!
```

but if

```bash
name="Lola!"
```

```
$ lova -r func.lova
LOVA: Translate: func.lova :: func.sh..
LOVA: Time compile: 451.588µs
Hello, Lova!
```

Nothing change

### Cases

Case - Consider a more convenient if/else 

`_` => analog else

```bash
proc init = "start"

case init do {
    "start":
        printn "Deploy start!"
        over
    "stop":
        printn "Deploy stop!"
        over
    _:
        printn "Error deploy!"
        over // analog break
}
```
### Defer
defer is something that must be executed at the end of the code, for example

```bash
// example defer
printn "Hello"
defer printn "World!"
```

```
$ lova -r defer.lova
LOVA: Translate: defer.lova :: defer.sh..
LOVA: Time compile: 215.54µs
Hello
World!
```

### Loop

1: While

While is basic loop in LOVA

```bash
// while

while true {
    print "\033[H\033[J"

    case frame do {
        "0":
            printn " /\\_/\\"
            printn "( o.o )"
            printn " > ^ < "
            over
        _:
            printn " /\\_/\\"
            printn "( -.- )"
            printn " > ^ < "
            over
    }

    sleep 1
    if frame == 0 {
        frame=1
    } else {
        frame=0
    }
}
```

while is one arguments

2: fdo

fdo loop = works as long as the conditions are false

```bash
// fdo

proc flag = true
fdo flag {
    printn "{COLOR_MAGENTA}[FDO] Flag loop check passed.{COLOR_NC}"
}
```

3: for in + arrays

for in = Used to iterate over arrays

```bash
// for in
// 1. Declare the array
arr [my_array] = ["Apple" "Banana" "Orange"]

// 2. Print all array elements (using [!] which parses to [@])
printn "All elements: [my_array][!]"

// 3. Append new elements
[my_array].append("Pear" "Kiwi")
printn "After append: [my_array][!]"

// 4. Delete element at index 1
[my_array].delete(1)
printn "After deleting index 1: [my_array][!]"

// 5. ranges for in
for el in [my_array] {
    printn "Fruits: $el"
}
```

## Imports
Import modules in Lova scripts

```bash
import "os"

if {IS_LINUX} == "true" {
    printn "You in Linux BTW!"
}
```

### shortcut

is_file (-f), is_dir (-d), is_exist (-e)
not (!), and (&&), or (||) 


# FAQ

- Why I Need LOVA Is Bash

`Yes, there is, but LOVA gives you more security with a convenient syntax, although knowing Bash doesn't hurt`

- What if I don't understand something?

`Go to discord I can often answer`

- What is SVS?

`This system works with the set -u flag, it prohibits making empty variables, it also prohibits making empty arguments $1, $2, an example of how to bypass if you need to declare a variable is already there`
