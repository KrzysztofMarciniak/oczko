package parser

import (
	"fmt"
	"strings"
	"unicode"
)

func Parse(source string) ([]Node, error) {
	tokens, err := lex(source)
	if err != nil {
		return nil, err
	}
	p := state{tokens: tokens}
	var nodes []Node
	for p.pos < len(p.tokens) {
		n, err := p.node()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

type state struct {
	tokens []string
	pos    int
}

func (p *state) node() (Node, error) {
	if p.pos >= len(p.tokens) {
		return nil, fmt.Errorf("unexpected end of input")
	}
	t := p.tokens[p.pos]
	p.pos++
	if t == ")" {
		return nil, fmt.Errorf("unexpected ')' ")
	}
	if t != "(" {
		return Atom{Value: t}, nil
	}
	var items []Node
	for {
		if p.pos >= len(p.tokens) {
			return nil, fmt.Errorf("unclosed '('")
		}
		if p.tokens[p.pos] == ")" {
			p.pos++
			break
		}
		n, err := p.node()
		if err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return List{Items: items}, nil
}
func lex(s string) ([]string, error) {
	var out []string
	for i := 0; i < len(s); {
		if unicode.IsSpace(rune(s[i])) {
			i++
			continue
		}
		if s[i] == ';' {
			for i < len(s) && s[i] != '\n' {
				i++
			}
			continue
		}
		if s[i] == '(' && i+1 < len(s) && s[i+1] == ';' {
			i += 2
			end := strings.Index(s[i:], "*)")
			if end < 0 {
				return nil, fmt.Errorf("unclosed block comment")
			}
			i += end + 2
			continue
		}
		if s[i] == '(' || s[i] == ')' {
			out = append(out, s[i:i+1])
			i++
			continue
		}
		if s[i] == '"' {
			i++
			var b strings.Builder
			closed := false
			for i < len(s) {
				if s[i] == '\\' && i+1 < len(s) {
					i++
					switch s[i] {
					case 'n':
						b.WriteByte('\n')
					case 't':
						b.WriteByte('\t')
					case '"':
						b.WriteByte('"')
					case '\\':
						b.WriteByte('\\')
					default:
						b.WriteByte(s[i])
					}
					i++
					continue
				}
				if s[i] == '"' {
					i++
					closed = true
					break
				}
				b.WriteByte(s[i])
				i++
			}
			if !closed {
				return nil, fmt.Errorf("unterminated string literal")
			}
			out = append(out, "\x00"+b.String())
			continue
		}
		start := i
		for i < len(s) && !unicode.IsSpace(rune(s[i])) && s[i] != '(' && s[i] != ')' {
			i++
		}
		out = append(out, s[start:i])
	}
	return out, nil
}
