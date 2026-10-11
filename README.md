# oczko

An esoteric programming language built around nested loops.

## Syntax

- `ENTER` — enter a new loop.
- `LOOP n` — repeat the current loop `n` additional times.
- `EXIT` — end the current loop.
- `(PRINT "text")` — print text.
- `(VAR name value)` — declare a local variable.
- `(GLOBAL name value)` — declare a global variable.
- `(SET name value)` — assign a variable.
- `(GET name)` — retrieve a variable's value.
- `(IF condition expression)` — execute an expression if the condition is true.
- `(; text)` - inline comment.

Expressions use Lisp-style prefix (Polish) notation, e.g. `(+ 1 2)`.

Local variables belong to their enclosing loop; global variables are accessible throughout the program.

Every loop must end with `EXIT`. An `IF` block must not introduce a separate loop ending; the enclosing loop has exactly one `EXIT`.

## Invalid Example

```text
(ENTER
    (IF (> (GET x) 0)
        (PRINT "Positive")
        EXIT ; invalid seperate loop ending 
    )
EXIT)
```
## Correct

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
## Factorial

Calculates \(5! = 120\) using a loop and `IF`.

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

