package php

import (
	"fmt"
	"strconv"
	"strings"
)

type token struct {
	kind string
	text string
	line int
}

type syntaxError struct {
	line    int
	message string
}

func (e *syntaxError) Error() string { return e.message }

func syntax(line int, format string, args ...any) error {
	return &syntaxError{line, fmt.Sprintf(format, args...)}
}

type lexer struct {
	source string
	pos    int
	line   int
}

func whitespace(c byte) bool { return strings.ContainsRune(" \t\r\n\v\f", rune(c)) }
func digit(c byte) bool      { return c >= '0' && c <= '9' }
func nameStart(c byte) bool {
	return c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= 0x80
}
func namePart(c byte) bool { return nameStart(c) || digit(c) }

func (l *lexer) take() byte {
	c := l.source[l.pos]
	l.pos++
	if c == '\n' || c == '\r' && (l.pos == len(l.source) || l.source[l.pos] != '\n') {
		l.line++
	}
	return c
}

func lex(source string) ([]token, error) {
	l := lexer{source: source, line: 1}
	if len(source) < 6 || !strings.EqualFold(source[:5], "<?php") || !whitespace(source[5]) {
		return nil, syntax(1, "unsupported syntax: file must begin with <?php followed by whitespace")
	}
	l.pos = 5
	var tokens []token
	for {
		t, err := l.next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
		if t.kind == "eof" {
			return tokens, nil
		}
	}
}

func (l *lexer) next() (token, error) {
	for l.pos < len(l.source) {
		rest := l.source[l.pos:]
		if whitespace(rest[0]) {
			l.take()
			continue
		}
		if strings.HasPrefix(rest, "#[") {
			return token{}, syntax(l.line, "unsupported syntax: attributes")
		}
		if strings.HasPrefix(rest, "//") || rest[0] == '#' {
			for l.pos < len(l.source) && l.source[l.pos] != '\n' && l.source[l.pos] != '\r' && !strings.HasPrefix(l.source[l.pos:], "?>") {
				l.take()
			}
			continue
		}
		if strings.HasPrefix(rest, "/*") {
			line := l.line
			l.take()
			l.take()
			for l.pos < len(l.source) && !strings.HasPrefix(l.source[l.pos:], "*/") {
				l.take()
			}
			if l.pos == len(l.source) {
				return token{}, syntax(line, "Unterminated comment starting line %d", line)
			}
			l.take()
			l.take()
			continue
		}
		break
	}
	t := token{line: l.line}
	if l.pos == len(l.source) {
		t.kind = "eof"
		return t, nil
	}
	rest := l.source[l.pos:]
	if strings.HasPrefix(rest, "?>") {
		l.take()
		l.take()
		if l.pos < len(l.source) && l.source[l.pos] == '\r' {
			l.take()
		}
		if l.pos < len(l.source) && l.source[l.pos] == '\n' {
			l.take()
		}
		if l.pos != len(l.source) {
			return t, syntax(l.line, "unsupported syntax: inline HTML or multiple PHP blocks")
		}
		t.kind = ";" // A final closing tag also terminates the statement.
		return t, nil
	}
	c := rest[0]
	if c == '\'' || c == '"' {
		t.kind = "string"
		value, err := l.quoted(c)
		t.text = value
		return t, err
	}
	if c == '$' {
		l.take()
		if l.pos == len(l.source) || !nameStart(l.source[l.pos]) {
			return t, syntax(t.line, "expected variable name after $")
		}
		start := l.pos
		for l.pos < len(l.source) && namePart(l.source[l.pos]) {
			l.take()
		}
		t.kind, t.text = "variable", l.source[start:l.pos]
		switch t.text {
		case "GLOBALS", "_SERVER", "_GET", "_POST", "_FILES", "_COOKIE", "_SESSION", "_REQUEST", "_ENV", "argc", "argv":
			return t, syntax(t.line, "unsupported predefined variable $%s", t.text)
		}
		return t, nil
	}
	if nameStart(c) {
		start := l.pos
		for l.pos < len(l.source) && namePart(l.source[l.pos]) {
			l.take()
		}
		t.kind, t.text = "name", l.source[start:l.pos]
		return t, nil
	}
	if digit(c) || c == '.' && len(rest) > 1 && digit(rest[1]) {
		start := l.pos
		for l.pos < len(l.source) && digit(l.source[l.pos]) {
			l.take()
		}
		integerEnd := l.pos
		if l.pos < len(l.source) && l.source[l.pos] == '.' {
			l.take()
			for l.pos < len(l.source) && digit(l.source[l.pos]) {
				l.take()
			}
		}
		if l.pos < len(l.source) && (l.source[l.pos] == 'e' || l.source[l.pos] == 'E') {
			l.take()
			if l.pos < len(l.source) && (l.source[l.pos] == '+' || l.source[l.pos] == '-') {
				l.take()
			}
			exponent := l.pos
			for l.pos < len(l.source) && digit(l.source[l.pos]) {
				l.take()
			}
			if exponent == l.pos {
				return t, syntax(t.line, "expected exponent digits")
			}
		}
		if l.pos < len(l.source) && nameStart(l.source[l.pos]) {
			return t, syntax(t.line, "unsupported numeric literal")
		}
		if integerEnd == l.pos && l.pos-start > 1 && c == '0' {
			return t, syntax(t.line, "unsupported numeric literal: octal")
		}
		t.kind, t.text = "number", l.source[start:l.pos]
		return t, nil
	}
	// Longest matches first, including !== before the unary ! token.
	for _, op := range []string{"===", "!==", "<=>", "++", "--", "+=", "-=", "*=", "/=", "%=", ".=", "&&", "||", "==", "!=", "<>", "<=", ">=", "=>", "->"} {
		if strings.HasPrefix(rest, op) {
			l.pos += len(op)
			t.kind, t.text = op, op
			return t, nil
		}
	}
	for _, op := range []string{"**", "<<", ">>"} {
		if strings.HasPrefix(rest, op) {
			return t, syntax(t.line, "unsupported operator %q", op)
		}
	}
	if strings.ContainsRune("=;,+-*/%.(){}[]!<>&", rune(c)) {
		l.take()
		t.kind, t.text = string(c), string(c)
		return t, nil
	}
	return t, syntax(t.line, "unsupported token %q", string(c))
}

func hex(c byte) bool { return digit(c) || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' }

func (l *lexer) quoted(quote byte) (string, error) {
	line := l.line
	start := l.pos
	l.take()
	var b strings.Builder
	for l.pos < len(l.source) {
		c := l.take()
		if c == quote {
			return b.String(), nil
		}
		if quote == '"' && c == '$' && l.pos < len(l.source) && (nameStart(l.source[l.pos]) || l.source[l.pos] == '{') {
			return "", syntax(l.line, "unsupported syntax: string interpolation")
		}
		if c != '\\' {
			b.WriteByte(c)
			continue
		}
		if l.pos == len(l.source) {
			break
		}
		next := l.take()
		if quote == '\'' {
			if next != '\\' && next != '\'' {
				b.WriteByte('\\')
			}
			b.WriteByte(next)
			continue
		}
		switch next {
		case 'n':
			b.WriteByte('\n')
		case 'r':
			b.WriteByte('\r')
		case 't':
			b.WriteByte('\t')
		case 'v':
			b.WriteByte('\v')
		case 'e':
			b.WriteByte(27)
		case 'f':
			b.WriteByte('\f')
		case '\\', '"', '$':
			b.WriteByte(next)
		case 'x':
			start := l.pos
			for l.pos < len(l.source) && l.pos-start < 2 && hex(l.source[l.pos]) {
				l.take()
			}
			if start == l.pos {
				b.WriteString("\\x")
			} else {
				n, _ := strconv.ParseUint(l.source[start:l.pos], 16, 8)
				b.WriteByte(byte(n))
			}
		case 'u':
			if l.pos < len(l.source) && l.source[l.pos] == '{' {
				return "", syntax(l.line, "unsupported syntax: Unicode codepoint escape")
			}
			b.WriteString("\\u")
		default:
			if next >= '0' && next <= '7' {
				n := int(next - '0')
				for count := 1; count < 3 && l.pos < len(l.source) && l.source[l.pos] >= '0' && l.source[l.pos] <= '7'; count++ {
					n = n*8 + int(l.take()-'0')
				}
				if n > 255 {
					return "", syntax(l.line, "unsupported syntax: overflowing octal escape")
				}
				b.WriteByte(byte(n))
			} else {
				b.WriteByte('\\')
				b.WriteByte(next)
			}
		}
	}
	if quote == '\'' {
		// PHP reports the raw, unclosed single-quoted token as string content.
		// Follow its diagnostic preview: first line, no surrounding quotes,
		// and 30 bytes plus an ellipsis when the content exceeds 33 bytes.
		content := l.source[start:]
		if newline := strings.IndexByte(content, '\n'); newline >= 0 {
			content = content[:newline]
		}
		content = content[1:]
		if strings.HasSuffix(content, "'") || strings.HasSuffix(content, "\"") {
			content = content[:len(content)-1]
		}
		if len(content) > 33 {
			content = content[:30] + "..."
		}
		return "", syntax(line, "syntax error, unexpected string content \"%s\"", content)
	}
	return "", syntax(line, "unterminated string")
}
