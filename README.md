# oczko

An esoteric programming language built around nested loops.

## Syntax

- `ENTER` — enter a new loop.
- `LOOP n` — repeat the current loop `n` additional times.
- `EXIT` — terminate the current loop.
- `(PRINT "text")` — print text.
- `(VAR name value)` — declare a local variable.
- `(GLOBAL name value)` — declare a global variable.
- `(SET name value)` — assign a new value to a variable.
- `(GET name)` — retrieve a variable's value.
- `(IF condition expression...)` — execute one or more expressions if the condition is true.
- `(; text)` — write a block comment.

Expressions use Lisp-style prefix (Polish) notation, e.g. `(+ 1 2)`.

Local variables belong to their enclosing loop scope. Global variables are accessible throughout the program.

Every loop must end with exactly one `EXIT`. An `IF` block does not introduce a separate loop scope or require its own loop terminator.

## Invalid Example

```text
(ENTER
    (IF (> (GET x) 0)
        (PRINT "Positive")
        EXIT ; invalid: separate loop terminator
    )
EXIT)
```

## Correct Example

```text
(ENTER
    (IF (> (GET x) 0)
        (PRINT "Positive")
    )
EXIT)
```

## Hello, World!

```text
(ENTER
    (PRINT "Hello, World!")
    (LOOP 1)
EXIT)
```

Output:

```text
Hello, World!
Hello, World!
```

`LOOP 1` executes the loop body twice: once initially and once more.

## Factorial

The following program calculates \(5! = 120\) using a nested loop and an `IF` statement.

```lisp
(ENTER
    (GLOBAL n 5)
    (GLOBAL result 1)

    (ENTER
        (IF (> (GET n) 1)
            (SET result (* (GET result) (GET n)))
            (SET n (- (GET n) 1))
        )
        (LOOP 4)
    EXIT)

    (PRINT (GET result))
EXIT)
```

Output:

```text
120
```

## Implementation

Oczko includes an interpreter written in Go.

### Requirements

- Go 1.22 or later
- Make (optional)

### Build and Test

```sh
make
make test
make clean
make install
```

### Running Programs

Run a source file:

```sh
./bin/oczko examples/hello.ocz
```

Or read source code from standard input:

```sh
cat examples/hello.ocz | ./bin/oczko
```

### Language Semantics

- **Loops:** `LOOP n` repeats the enclosing loop `n` additional times, so its body executes `n+1` times. `LOOP` must appear at the loop's top level.
- **Loop termination:** Every `ENTER` requires exactly one matching `EXIT`. `EXIT` is structural syntax, not an executable statement.
- **Conditionals:** `IF` accepts a condition followed by one or more expressions. It does not create a new loop.
- **Local variables:** `VAR` declares a variable in the current loop scope.
- **Global variables:** `GLOBAL` declares a program-wide variable.
- **Assignment:** `SET` updates the nearest matching local variable or, if none exists, a matching global variable.
- **Variable lookup:** `GET` searches local scopes from nearest to farthest, then checks global variables.
- **Values:** The interpreter supports integers, strings, and booleans. Falsey values are `false`, `0`, empty strings, and `nil`.

### Supported Operators

| Category | Operators |
|---|---|
| Arithmetic | `+`, `-`, `*`, `/`, `%` |
| Comparison | `=`, `!=`, `<`, `<=`, `>`, `>=` |
| Logical | `and`, `or`, `not` |

All operators use prefix notation. For example:

```lisp
(+ 1 2)
(* (GET n) 2)
(> (GET n) 1)
```

### Comments

Oczko supports semicolon-to-end-of-line comments and block comments:

```text
; This is a line comment.

(; This is a block comment *)
```
