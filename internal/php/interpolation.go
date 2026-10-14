package php

import (
	"strconv"
	"strings"
)

type stringPart struct {
	text   string
	tokens []token
	line   int
}
type interpolatedString struct {
	parts []expr
	line  int
}

func (interpolatedString) expression() {}

func (p *parser) interpolated(t token) (expr, error) {
	n := interpolatedString{line: t.line}
	for _, part := range t.parts {
		if part.tokens == nil {
			n.parts = append(n.parts, literal{value{kind: stringKind, s: part.text}, part.line})
			continue
		}
		child := *p
		child.tokens, child.pos = part.tokens, 0
		x, err := child.expression(0)
		if err != nil {
			return nil, err
		}
		if child.current().kind != "eof" {
			return nil, syntax(part.line, "unsupported syntax: string interpolation expression")
		}
		if !interpolationAccess(x) {
			return nil, syntax(part.line, "unsupported syntax: string interpolation expression")
		}
		n.parts = append(n.parts, x)
	}
	return n, nil
}

func interpolationAccess(x expr) bool {
	switch n := x.(type) {
	case variable:
		return true
	case offsetExpr:
		return n.key != nil && interpolationAccess(n.base)
	case propertyExpr:
		return interpolationAccess(n.base)
	}
	return false
}

func (e *evaluator) evalInterpolated(n interpolatedString) (value, error) {
	var b strings.Builder
	for _, part := range n.parts {
		v, err := e.eval(part)
		if err != nil {
			return value{}, err
		}
		s, err := e.toText(v, expressionLine(part))
		if err != nil {
			return value{}, err
		}
		b.WriteString(s)
	}
	return value{kind: stringKind, s: b.String()}, nil
}

// Lex string contents separately from PHP expressions. Escapes are consumed
// before detecting '$'; interpolations become tokens/AST, never replacements.
func (l *lexer) doubleQuoted() (token, error) {
	line := l.line
	l.take()
	return l.stringContents('"', line)
}

func (l *lexer) stringContents(terminator byte, line int) (token, error) {
	t := token{kind: "interpolated", line: line}
	var b strings.Builder
	flush := func() {
		if b.Len() != 0 {
			t.parts = append(t.parts, stringPart{text: b.String(), line: l.line})
			b.Reset()
		}
	}
	for l.pos < len(l.source) {
		c := l.source[l.pos]
		if terminator != 0 && c == terminator {
			l.take()
			flush()
			return collapseString(t), nil
		}
		if c == '\\' {
			start := l.pos
			l.take()
			if l.pos == len(l.source) {
				b.WriteByte('\\')
				break
			}
			next := l.take()
			switch next {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case 'v':
				b.WriteByte('\v')
			case 'f':
				b.WriteByte('\f')
			case 'e':
				b.WriteByte(27)
			case '\\', '$':
				b.WriteByte(next)
			case '"':
				if terminator == '"' {
					b.WriteByte(next)
				} else {
					b.WriteString(`\"`)
				}
			case 'x':
				from := l.pos
				for l.pos < len(l.source) && l.pos-from < 2 && hex(l.source[l.pos]) {
					l.take()
				}
				if from == l.pos {
					b.WriteString(`\x`)
				} else {
					n, _ := strconv.ParseUint(l.source[from:l.pos], 16, 8)
					b.WriteByte(byte(n))
				}
			case 'u':
				if l.pos < len(l.source) && l.source[l.pos] == '{' {
					return t, syntax(l.line, "unsupported syntax: Unicode codepoint escape")
				}
				b.WriteString(`\u`)
			default:
				if next >= '0' && next <= '7' {
					n := int(next - '0')
					for count := 1; count < 3 && l.pos < len(l.source) && l.source[l.pos] >= '0' && l.source[l.pos] <= '7'; count++ {
						n = n*8 + int(l.take()-'0')
					}
					if n > 255 {
						return t, syntax(l.line, "unsupported syntax: overflowing octal escape")
					}
					b.WriteByte(byte(n))
				} else {
					b.WriteString(l.source[start:l.pos])
				}
			}
			continue
		}
		if strings.HasPrefix(l.source[l.pos:], "${") {
			return t, syntax(l.line, "unsupported syntax: deprecated ${} interpolation")
		}
		braced := strings.HasPrefix(l.source[l.pos:], "{$")
		if braced || c == '$' && l.pos+1 < len(l.source) && nameStart(l.source[l.pos+1]) {
			flush()
			part := stringPart{line: l.line}
			var err error
			if braced {
				part.tokens, err = l.bracedInterpolation()
			} else {
				part.tokens, err = l.simpleInterpolation()
			}
			if err != nil {
				return t, err
			}
			t.parts = append(t.parts, part)
			continue
		}
		b.WriteByte(l.take())
	}
	if terminator != 0 {
		return t, syntax(line, "unterminated string")
	}
	flush()
	return collapseString(t), nil
}

func collapseString(t token) token {
	var text strings.Builder
	for _, p := range t.parts {
		if p.tokens != nil {
			return t
		}
		text.WriteString(p.text)
	}
	t.kind, t.text, t.parts = "string", text.String(), nil
	return t
}

func (l *lexer) bracedInterpolation() ([]token, error) {
	l.take() // '{'
	var tokens []token
	for {
		t, err := l.next()
		if err != nil {
			return nil, err
		}
		if t.kind == "}" {
			tokens = append(tokens, token{kind: "eof", line: t.line})
			return tokens, nil
		}
		if t.kind == "eof" {
			return nil, syntax(t.line, "syntax error, unexpected end of file, expecting \"}\"")
		}
		if t.kind == "{" || t.kind == "interpolated" {
			return nil, syntax(t.line, "unsupported syntax: nested interpolation")
		}
		tokens = append(tokens, t)
	}
}

func (l *lexer) simpleInterpolation() ([]token, error) {
	variable, err := l.next()
	if err != nil {
		return nil, err
	}
	tokens := []token{variable}
	if l.pos < len(l.source) && l.source[l.pos] == '[' {
		line := l.line
		l.take()
		tokens = append(tokens, token{kind: "[", line: line})
		if l.pos < len(l.source) && l.source[l.pos] == '$' {
			v, err := l.next()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, v)
		} else {
			start := l.pos
			if l.pos < len(l.source) && l.source[l.pos] == '-' {
				l.take()
			}
			for l.pos < len(l.source) && namePart(l.source[l.pos]) {
				l.take()
			}
			key := l.source[start:l.pos]
			if key == "" {
				return nil, syntax(line, "unsupported syntax: simple interpolation offset")
			}
			if _, err := strconv.ParseInt(key, 10, 64); err == nil {
				if key[0] == '-' {
					tokens = append(tokens, token{kind: "-", text: "-", line: line})
					key = key[1:]
				}
				tokens = append(tokens, token{kind: "number", text: key, line: line})
			} else {
				tokens = append(tokens, token{kind: "string", text: key, line: line})
			}
		}
		if l.pos >= len(l.source) || l.source[l.pos] != ']' {
			return nil, syntax(line, "unsupported syntax: simple interpolation offset")
		}
		l.take()
		tokens = append(tokens, token{kind: "]", line: line})
	} else if strings.HasPrefix(l.source[l.pos:], "->") && l.pos+2 < len(l.source) && nameStart(l.source[l.pos+2]) {
		line := l.line
		l.take()
		l.take()
		start := l.pos
		for l.pos < len(l.source) && namePart(l.source[l.pos]) {
			l.take()
		}
		tokens = append(tokens, token{kind: "->", line: line}, token{kind: "name", text: l.source[start:l.pos], line: line})
	}
	tokens = append(tokens, token{kind: "eof", line: l.line})
	return tokens, nil
}

func (l *lexer) heredoc() (token, error) {
	line := l.line
	l.pos += 3
	for l.pos < len(l.source) && (l.source[l.pos] == ' ' || l.source[l.pos] == '\t') {
		l.take()
	}
	quote := byte(0)
	if l.pos < len(l.source) && (l.source[l.pos] == '\'' || l.source[l.pos] == '"') {
		quote = l.take()
	}
	start := l.pos
	if l.pos >= len(l.source) || !nameStart(l.source[l.pos]) {
		return token{}, syntax(line, "invalid heredoc label")
	}
	for l.pos < len(l.source) && namePart(l.source[l.pos]) {
		l.take()
	}
	label := l.source[start:l.pos]
	if quote != 0 {
		if l.pos >= len(l.source) || l.take() != quote {
			return token{}, syntax(line, "invalid heredoc label")
		}
	}
	if l.pos < len(l.source) && l.source[l.pos] == '\r' {
		l.take()
	}
	if l.pos >= len(l.source) || l.take() != '\n' {
		return token{}, syntax(line, "invalid heredoc opening line")
	}
	bodyLine, bodyStart := l.line, l.pos
	for l.pos < len(l.source) {
		lineStart := l.pos
		for l.pos < len(l.source) && (l.source[l.pos] == ' ' || l.source[l.pos] == '\t') {
			l.take()
		}
		indent := l.source[lineStart:l.pos]
		if strings.HasPrefix(l.source[l.pos:], label) && (l.pos+len(label) == len(l.source) || !namePart(l.source[l.pos+len(label)])) {
			body := l.source[bodyStart:lineStart]
			if strings.HasSuffix(body, "\r\n") {
				body = body[:len(body)-2]
			} else if strings.HasSuffix(body, "\n") {
				body = body[:len(body)-1]
			}
			l.pos += len(label)
			if strings.Contains(indent, " ") && strings.Contains(indent, "\t") {
				return token{}, syntax(l.line, "Invalid indentation - tabs and spaces cannot be mixed")
			}
			lines := strings.Split(body, "\n")
			for i, text := range lines {
				count := 0
				for count < len(text) && count < len(indent) && (text[count] == ' ' || text[count] == '\t') {
					if text[count] != indent[count] {
						return token{}, syntax(bodyLine+i, "Invalid indentation - tabs and spaces cannot be mixed")
					}
					count++
				}
				if count < len(indent) && strings.TrimSpace(text) != "" {
					return token{}, syntax(bodyLine+i, "Invalid body indentation level (expecting an indentation level of at least %d)", len(indent))
				}
				lines[i] = text[count:]
			}
			body = strings.Join(lines, "\n")
			if quote == '\'' {
				return token{kind: "string", text: body, line: line}, nil
			}
			child := lexer{source: body, line: bodyLine}
			return child.stringContents(0, line)
		}
		for l.pos < len(l.source) && l.take() != '\n' {
		}
	}
	return token{}, syntax(l.line, "syntax error, unexpected end of file, expecting heredoc end")
}
