package interpreter

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"oczko/src/impl/parser"
	"oczko/src/interface/runner"
)

type Interpreter struct {
	out     io.Writer
	globals map[string]any
}

var _ runner.Runner = (*Interpreter)(nil)

func New(out io.Writer) *Interpreter { return &Interpreter{out: out, globals: map[string]any{}} }

type scope struct {
	vars   map[string]any
	parent *scope
}

func newScope(parent *scope) *scope { return &scope{vars: map[string]any{}, parent: parent} }
func (s *scope) get(name string) (any, bool) {
	for p := s; p != nil; p = p.parent {
		v, ok := p.vars[name]
		if ok {
			return v, true
		}
	}
	return nil, false
}
func (i *Interpreter) Run(source string) error {
	nodes, err := parser.Parse(source)
	if err != nil {
		return err
	}
	root := newScope(nil)
	for _, n := range nodes {
		if _, err = i.exec(n, root); err != nil {
			return err
		}
	}
	return nil
}
func atom(n parser.Node) (string, bool)        { a, ok := n.(parser.Atom); return a.Value, ok }
func list(n parser.Node) ([]parser.Node, bool) { a, ok := n.(parser.List); return a.Items, ok }
func (i *Interpreter) exec(n parser.Node, s *scope) (any, error) {
	if a, ok := atom(n); ok {
		return i.literal(a)
	}
	xs, ok := list(n)
	if !ok || len(xs) == 0 {
		return nil, nil
	}
	head, ok := atom(xs[0])
	if !ok {
		return nil, fmt.Errorf("expression name must be an atom")
	}
	args := xs[1:]
	switch strings.ToUpper(head) {
	case "ENTER":
		return nil, i.execLoop(args, s)
	case "EXIT":
		return nil, fmt.Errorf("EXIT without matching ENTER")
	case "LOOP":
		return nil, fmt.Errorf("LOOP is only valid directly inside ENTER")
	case "PRINT":
		if len(args) != 1 {
			return nil, fmt.Errorf("PRINT expects one argument")
		}
		v, e := i.eval(args[0], s)
		if e != nil {
			return nil, e
		}
		_, e = fmt.Fprintln(i.out, stringify(v))
		return nil, e
	case "VAR", "GLOBAL":
		if len(args) != 2 {
			return nil, fmt.Errorf("%s expects name and value", head)
		}
		name, e := nameArg(args[0])
		if e != nil {
			return nil, e
		}
		v, e := i.eval(args[1], s)
		if e != nil {
			return nil, e
		}
		if strings.EqualFold(head, "GLOBAL") {
			i.globals[name] = v
		} else {
			s.vars[name] = v
		}
		return v, nil
	case "SET":
		if len(args) != 2 {
			return nil, fmt.Errorf("SET expects name and value")
		}
		name, e := nameArg(args[0])
		if e != nil {
			return nil, e
		}
		v, e := i.eval(args[1], s)
		if e != nil {
			return nil, e
		}
		for p := s; p != nil; p = p.parent {
			if _, found := p.vars[name]; found {
				p.vars[name] = v
				return v, nil
			}
		}
		if _, found := i.globals[name]; found {
			i.globals[name] = v
			return v, nil
		}
		return nil, fmt.Errorf("SET: undefined variable %q", name)
	case "GET":
		if len(args) != 1 {
			return nil, fmt.Errorf("GET expects one argument")
		}
		name, e := nameArg(args[0])
		if e != nil {
			return nil, e
		}
		if v, found := s.get(name); found {
			return v, nil
		}
		if v, found := i.globals[name]; found {
			return v, nil
		}
		return nil, fmt.Errorf("GET: undefined variable %q", name)
	case "IF":
		if len(args) < 2 {
			return nil, fmt.Errorf("IF expects a condition and at least one expression")
		}
		cond, e := i.eval(args[0], s)
		if e != nil {
			return nil, e
		}
		if truthy(cond) {
			var v any
			for _, a := range args[1:] {
				v, e = i.exec(a, s)
				if e != nil {
					return nil, e
				}
			}
			return v, nil
		}
		return nil, nil
	default:
		return i.eval(n, s)
	}
}
func (i *Interpreter) execLoop(args []parser.Node, parent *scope) error {
	body := make([]parser.Node, 0, len(args))
	additional := int64(0)
	exits := 0
	for _, n := range args {
		if a, ok := atom(n); ok && strings.EqualFold(a, "EXIT") {
			exits++
			continue
		}
		if xs, ok := list(n); ok && len(xs) > 0 {
			if h, ok := atom(xs[0]); ok && strings.EqualFold(h, "LOOP") {
				if len(xs) != 2 {
					return fmt.Errorf("LOOP expects one integer")
				}
				v, e := i.eval(xs[1], parent)
				if e != nil {
					return e
				}
				count, e := integer(v)
				if e != nil {
					return fmt.Errorf("LOOP: %w", e)
				}
				if count < 0 {
					return fmt.Errorf("LOOP count cannot be negative")
				}
				additional = count
				continue
			}
		}
		body = append(body, n)
	}
	if exits != 1 {
		return fmt.Errorf("ENTER loop must contain exactly one EXIT (found %d)", exits)
	}
	for iteration := int64(0); iteration <= additional; iteration++ {
		local := newScope(parent)
		for _, n := range body {
			if _, err := i.exec(n, local); err != nil {
				return err
			}
		}
	}
	return nil
}
func (i *Interpreter) eval(n parser.Node, s *scope) (any, error) {
	if a, ok := atom(n); ok {
		return i.literal(a)
	}
	xs, ok := list(n)
	if !ok || len(xs) == 0 {
		return nil, nil
	}
	h, ok := atom(xs[0])
	if !ok {
		return nil, fmt.Errorf("invalid expression")
	}
	if strings.EqualFold(h, "GET") {
		return i.exec(n, s)
	}
	vals := make([]any, 0, len(xs)-1)
	for _, x := range xs[1:] {
		v, e := i.eval(x, s)
		if e != nil {
			return nil, e
		}
		vals = append(vals, v)
	}
	nums := func() (int64, int64, error) {
		if len(vals) != 2 {
			return 0, 0, fmt.Errorf("%s expects two arguments", h)
		}
		a, e := integer(vals[0])
		if e != nil {
			return 0, 0, e
		}
		b, e := integer(vals[1])
		return a, b, e
	}
	switch strings.ToLower(h) {
	case "+", "-", "*", "/", "%":
		a, b, e := nums()
		if e != nil {
			return nil, e
		}
		switch h {
		case "+":
			return a + b, nil
		case "-":
			return a - b, nil
		case "*":
			return a * b, nil
		case "/":
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return a / b, nil
		default:
			if b == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			return a % b, nil
		}
	case ">", ">=", "<", "<=", "=", "!=":
		if len(vals) != 2 {
			return nil, fmt.Errorf("%s expects two arguments", h)
		}
		a, e := integer(vals[0])
		if e == nil {
			b, e := integer(vals[1])
			if e != nil {
				return nil, e
			}
			switch h {
			case ">":
				return a > b, nil
			case ">=":
				return a >= b, nil
			case "<":
				return a < b, nil
			case "<=":
				return a <= b, nil
			case "=":
				return a == b, nil
			case "!=":
				return a != b, nil
			}
		}
		if h == "=" {
			return vals[0] == vals[1], nil
		}
		if h == "!=" {
			return vals[0] != vals[1], nil
		}
		return nil, fmt.Errorf("%s requires numeric arguments", h)
	case "and":
		for _, v := range vals {
			if !truthy(v) {
				return false, nil
			}
		}
		return true, nil
	case "or":
		for _, v := range vals {
			if truthy(v) {
				return true, nil
			}
		}
		return false, nil
	case "not":
		if len(vals) != 1 {
			return nil, fmt.Errorf("not expects one argument")
		}
		return !truthy(vals[0]), nil
	default:
		return nil, fmt.Errorf("unknown expression %q", h)
	}
}
func (i *Interpreter) literal(a string) (any, error) {
	if strings.HasPrefix(a, "\x00") {
		return strings.TrimPrefix(a, "\x00"), nil
	}
	if a == "true" {
		return true, nil
	}
	if a == "false" {
		return false, nil
	}
	if n, e := strconv.ParseInt(a, 10, 64); e == nil {
		return n, nil
	}
	return a, nil
}
func nameArg(n parser.Node) (string, error) {
	a, ok := atom(n)
	if !ok {
		return "", fmt.Errorf("variable name must be an atom")
	}
	return a, nil
}
func integer(v any) (int64, error) {
	switch n := v.(type) {
	case int64:
		return n, nil
	case int:
		return int64(n), nil
	default:
		return 0, fmt.Errorf("expected integer, got %T (%v)", v, v)
	}
}
func truthy(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case int64:
		return x != 0
	case int:
		return x != 0
	case string:
		return x != ""
	}
	return true
}
func stringify(v any) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprint(v)
}
