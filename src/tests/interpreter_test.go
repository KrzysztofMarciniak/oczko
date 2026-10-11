package tests

import (
	"bytes"
	"oczko/src/impl/interpreter"
	"testing"
)

func run(t *testing.T, source string) string {
	t.Helper()
	var out bytes.Buffer
	if err := interpreter.New(&out).Run(source); err != nil {
		t.Fatal(err)
	}
	return out.String()
}
func TestHelloWorld(t *testing.T) {
	got := run(t, `(ENTER (PRINT "Hello, World!") (LOOP 1) EXIT)`)
	want := "Hello, World!\nHello, World!\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
func TestFactorial(t *testing.T) {
	got := run(t, `(ENTER (GLOBAL n 5) (GLOBAL result 1) (ENTER (IF (> (GET n) 1) (SET result (* (GET result) (GET n))) (SET n (- (GET n) 1))) (LOOP 4) EXIT) (PRINT (GET result)) EXIT)`)
	if got != "120\n" {
		t.Fatalf("got %q", got)
	}
}
func TestLocalScopeAndIf(t *testing.T) {
	got := run(t, `(ENTER (VAR x 2) (IF (> (GET x) 0) (PRINT "yes")) EXIT)`)
	if got != "yes\n" {
		t.Fatalf("got %q", got)
	}
}
func TestComments(t *testing.T) {
	got := run(t, `(; block comment *) (ENTER (PRINT "ok") ; trailing comment
 EXIT)`)
	if got != "ok\n" {
		t.Fatalf("got %q", got)
	}
}
