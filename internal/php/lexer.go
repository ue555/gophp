package php

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type token struct {
	kind  string
	text  string
	line  int
	start int
	parts []stringPart
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
	source  string
	pos     int
	line    int
	phpMode bool
	pending []token
}

func whitespace(c byte) bool { return strings.ContainsRune(" \t\r\n\v\f", rune(c)) }
func digit(c byte) bool      { return c >= '0' && c <= '9' }
func (l *lexer) numericDigit(chars string, start int) bool {
	if l.pos >= len(l.source) {
		return false
	}
	if strings.ContainsRune(chars, rune(l.source[l.pos])) {
		return true
	}
	return l.source[l.pos] == '_' && l.pos > start && l.pos+1 < len(l.source) && strings.ContainsRune(chars, rune(l.source[l.pos+1]))
}
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
	if strings.HasPrefix(source, "#!") {
		if end := strings.IndexByte(source, '\n'); end >= 0 {
			l.pos = end + 1
			l.line = 2
		}
	}
	if len(source)-l.pos >= 6 && strings.EqualFold(source[l.pos:l.pos+5], "<?php") && whitespace(source[l.pos+5]) {
		l.pos += 5
		l.phpMode = true
	}
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
	if len(l.pending) > 0 {
		t := l.pending[0]
		l.pending = l.pending[1:]
		return t, nil
	}
	if !l.phpMode {
		start, line := l.pos, l.line
		for l.pos < len(l.source) {
			rest := l.source[l.pos:]
			if len(rest) >= 6 && strings.EqualFold(rest[:5], "<?php") && whitespace(rest[5]) || strings.HasPrefix(rest, "<?=") {
				break
			}
			l.take()
		}
		if l.pos > start {
			l.pending = []token{{kind: "string", text: l.source[start:l.pos], line: line, start: start}, {kind: ";", line: line, start: start}}
			return token{kind: "name", text: "echo", line: line, start: start}, nil
		}
		if l.pos == len(l.source) {
			return token{kind: "eof", line: l.line, start: l.pos}, nil
		}
		if strings.HasPrefix(l.source[l.pos:], "<?=") {
			l.pos += 3
			l.phpMode = true
			return token{kind: "name", text: "echo", line: l.line, start: l.pos - 3}, nil
		}
		l.pos += 5
		l.phpMode = true
	}
	for l.pos < len(l.source) {
		rest := l.source[l.pos:]
		if whitespace(rest[0]) {
			l.take()
			continue
		}
		if strings.HasPrefix(rest, "#[") {
			break
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
	t := token{line: l.line, start: l.pos}
	if l.pos == len(l.source) {
		t.kind = "eof"
		return t, nil
	}
	rest := l.source[l.pos:]
	if strings.HasPrefix(rest, "#[") {
		l.take()
		l.take()
		t.kind = "#["
		return t, nil
	}
	if strings.HasPrefix(rest, "?>") {
		l.take()
		l.take()
		if l.pos < len(l.source) && l.source[l.pos] == '\r' {
			l.take()
		}
		if l.pos < len(l.source) && l.source[l.pos] == '\n' {
			l.take()
		}
		l.phpMode = false
		t.kind = ";" // A final closing tag also terminates the statement.
		return t, nil
	}
	c := rest[0]
	if strings.HasPrefix(rest, "<<<") {
		return l.heredoc()
	}
	if c == '"' {
		return l.doubleQuoted()
	}
	if c == '\'' || c == '"' {
		t.kind = "string"
		value, err := l.quoted(c)
		t.text = value
		return t, err
	}
	if c == '$' {
		l.take()
		if l.pos < len(l.source) && (l.source[l.pos] == '{' || l.source[l.pos] == '$') {
			t.kind = "$"
			return t, nil
		}
		if l.pos == len(l.source) || !nameStart(l.source[l.pos]) {
			return t, syntax(t.line, "expected variable name after $")
		}
		start := l.pos
		for l.pos < len(l.source) && namePart(l.source[l.pos]) {
			l.take()
		}
		t.kind, t.text = "variable", l.source[start:l.pos]
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
		if len(rest) > 2 && rest[0] == '0' && (rest[1] == 'x' || rest[1] == 'X') {
			l.take()
			l.take()
			digits := l.pos
			for l.numericDigit("0123456789abcdefABCDEF", digits) {
				l.take()
			}
			if l.pos == digits || l.pos < len(l.source) && namePart(l.source[l.pos]) {
				return t, syntax(t.line, "unsupported numeric literal")
			}
			t.kind, t.text = "number", l.source[start:l.pos]
			return t, nil
		}
		if len(rest) > 2 && rest[0] == '0' && (rest[1] == 'o' || rest[1] == 'O' || rest[1] == 'b' || rest[1] == 'B') {
			base := "01234567"
			if rest[1] == 'b' || rest[1] == 'B' {
				base = "01"
			}
			l.take()
			l.take()
			digits := l.pos
			for l.numericDigit(base, digits) {
				l.take()
			}
			if l.pos == digits || l.pos < len(l.source) && namePart(l.source[l.pos]) {
				return t, syntax(t.line, "unsupported numeric literal")
			}
			t.kind, t.text = "number", l.source[start:l.pos]
			return t, nil
		}
		for l.numericDigit("0123456789", start) {
			l.take()
		}
		integerEnd := l.pos
		if l.pos < len(l.source) && l.source[l.pos] == '.' {
			l.take()
			for l.numericDigit("0123456789", start) {
				l.take()
			}
		}
		if l.pos < len(l.source) && (l.source[l.pos] == 'e' || l.source[l.pos] == 'E') {
			l.take()
			if l.pos < len(l.source) && (l.source[l.pos] == '+' || l.source[l.pos] == '-') {
				l.take()
			}
			exponent := l.pos
			for l.numericDigit("0123456789", start) {
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
			for _, digit := range l.source[start+1 : l.pos] {
				if digit == '_' {
					continue
				}
				if digit < '0' || digit > '7' {
					return t, syntax(t.line, "Invalid numeric literal")
				}
			}
		}
		t.kind, t.text = "number", l.source[start:l.pos]
		return t, nil
	}
	// Longest matches first, including !== before the unary ! token.
	for _, op := range []string{"??=", "?->", "<<=", ">>=", "===", "!==", "<=>", "...", "??", "<<", ">>", "&=", "|=", "^=", "::", "++", "--", "+=", "-=", "*=", "/=", "%=", ".=", "**", "&&", "||", "==", "!=", "<>", "<=", ">=", "=>", "->"} {
		if strings.HasPrefix(rest, op) {
			l.pos += len(op)
			t.kind, t.text = op, op
			return t, nil
		}
	}
	if strings.ContainsRune("=;,+-*/%.(){}[]!<>&|^~?:@\\", rune(c)) {
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
				l.take()
				start := l.pos
				for l.pos < len(l.source) && hex(l.source[l.pos]) {
					l.take()
				}
				if l.pos == start || l.pos >= len(l.source) || l.source[l.pos] != '}' {
					return "", syntax(l.line, "Invalid UTF-8 codepoint escape sequence")
				}
				code, err := strconv.ParseUint(l.source[start:l.pos], 16, 32)
				l.take()
				if err != nil || code > utf8.MaxRune || code >= 0xD800 && code <= 0xDFFF {
					return "", syntax(l.line, "Invalid UTF-8 codepoint escape sequence")
				}
				b.WriteRune(rune(code))
				break
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
