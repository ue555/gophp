package php

import (
	"regexp"
	"strconv"
	"strings"
)

type phpLexeme struct {
	name, text string
	line       int
}
type phpTokenizer struct {
	source     string
	pos, line  int
	php, parse bool
	version    Version
	tokens     []phpLexeme
}

var tokenNumber = regexp.MustCompile(`^(?:0[xX][0-9a-fA-F](?:_?[0-9a-fA-F])*|0[bB][01](?:_?[01])*|0[oO][0-7](?:_?[0-7])*|(?:[0-9](?:_?[0-9])*(?:\.(?:[0-9](?:_?[0-9])*)?)?|\.[0-9](?:_?[0-9])*)(?:[eE][+-]?[0-9](?:_?[0-9])*)?)`)
var tokenHeredoc = regexp.MustCompile(`^<<<[ \t]*(['"]?)([A-Za-z_\x80-\xff][A-Za-z0-9_\x80-\xff]*)['"]?\r?\n`)
var phpOperators = [][2]string{{"?->", "T_NULLSAFE_OBJECT_OPERATOR"}, {"===", "T_IS_IDENTICAL"}, {"!==", "T_IS_NOT_IDENTICAL"}, {"<=>", "T_SPACESHIP"}, {"<<=", "T_SL_EQUAL"}, {">>=", "T_SR_EQUAL"}, {"**=", "T_POW_EQUAL"}, {"??=", "T_COALESCE_EQUAL"}, {"...", "T_ELLIPSIS"}, {"->", "T_OBJECT_OPERATOR"}, {"::", "T_DOUBLE_COLON"}, {"=>", "T_DOUBLE_ARROW"}, {"==", "T_IS_EQUAL"}, {"!=", "T_IS_NOT_EQUAL"}, {"<>", "T_IS_NOT_EQUAL"}, {"<=", "T_IS_SMALLER_OR_EQUAL"}, {">=", "T_IS_GREATER_OR_EQUAL"}, {"++", "T_INC"}, {"--", "T_DEC"}, {"+=", "T_PLUS_EQUAL"}, {"-=", "T_MINUS_EQUAL"}, {"*=", "T_MUL_EQUAL"}, {"/=", "T_DIV_EQUAL"}, {".=", "T_CONCAT_EQUAL"}, {"%=", "T_MOD_EQUAL"}, {"&=", "T_AND_EQUAL"}, {"|=", "T_OR_EQUAL"}, {"^=", "T_XOR_EQUAL"}, {"&&", "T_BOOLEAN_AND"}, {"||", "T_BOOLEAN_OR"}, {"<<", "T_SL"}, {">>", "T_SR"}, {"**", "T_POW"}, {"??", "T_COALESCE"}, {"#[", "T_ATTRIBUTE"}, {"|>", "T_PIPE"}}
var phpKeywords = map[string]string{}

func init() {
	for _, w := range strings.Fields("include include_once eval require require_once print yield instanceof new clone exit if elseif else endif echo do while endwhile for endfor foreach endforeach declare enddeclare as switch endswitch case default match break continue goto function fn const return try catch finally throw use insteadof global static abstract final private protected public readonly var unset isset empty halt_compiler class trait interface enum extends implements namespace list array callable") {
		phpKeywords[w] = "T_" + strings.ToUpper(w)
	}
	for a, b := range map[string]string{"die": "T_EXIT", "or": "T_LOGICAL_OR", "and": "T_LOGICAL_AND", "xor": "T_LOGICAL_XOR", "__halt_compiler": "T_HALT_COMPILER", "__line__": "T_LINE", "__file__": "T_FILE", "__dir__": "T_DIR", "__class__": "T_CLASS_C", "__trait__": "T_TRAIT_C", "__method__": "T_METHOD_C", "__function__": "T_FUNC_C", "__namespace__": "T_NS_C", "__property__": "T_PROPERTY_C"} {
		phpKeywords[a] = b
	}
}
func (t *phpTokenizer) emit(name string, end int) {
	start, line := t.pos, t.line
	for t.pos < end {
		c := t.source[t.pos]
		t.pos++
		if c == '\n' || c == '\r' && (t.pos == len(t.source) || t.source[t.pos] != '\n') {
			t.line++
		}
	}
	t.tokens = append(t.tokens, phpLexeme{name, t.source[start:end], line})
}
func (t *phpTokenizer) previous() string {
	for i := len(t.tokens) - 1; i >= 0; i-- {
		n := t.tokens[i].name
		if n != "T_WHITESPACE" && n != "T_COMMENT" && n != "T_DOC_COMMENT" {
			return n
		}
	}
	return ""
}
func (t *phpTokenizer) scan() error {
	for t.pos < len(t.source) {
		rest := t.source[t.pos:]
		if !t.php {
			if strings.HasPrefix(rest, "<?=") {
				t.emit("T_OPEN_TAG_WITH_ECHO", t.pos+3)
				t.php = true
				continue
			}
			if len(rest) >= 5 && strings.EqualFold(rest[:5], "<?php") && (len(rest) == 5 || whitespace(rest[5])) {
				end := t.pos + 5
				if end < len(t.source) {
					if strings.HasPrefix(t.source[end:], "\r\n") {
						end += 2
					} else {
						end++
					}
				}
				t.emit("T_OPEN_TAG", end)
				t.php = true
				continue
			}
			end := t.pos + 1
			for end < len(t.source) {
				s := t.source[end:]
				if strings.HasPrefix(s, "<?=") || len(s) >= 5 && strings.EqualFold(s[:5], "<?php") && (len(s) == 5 || whitespace(s[5])) {
					break
				}
				end++
			}
			t.emit("T_INLINE_HTML", end)
			continue
		}
		if err := t.one(); err != nil {
			return err
		}
	}
	return nil
}
func (t *phpTokenizer) one() error {
	rest := t.source[t.pos:]
	c := rest[0]
	start := t.pos
	if whitespace(c) {
		end := start + 1
		for end < len(t.source) && whitespace(t.source[end]) {
			end++
		}
		t.emit("T_WHITESPACE", end)
		return nil
	}
	if strings.HasPrefix(rest, "?>") {
		end := start + 2
		if strings.HasPrefix(t.source[end:], "\r\n") {
			end += 2
		} else if end < len(t.source) && (t.source[end] == '\n' || t.source[end] == '\r') {
			end++
		}
		t.emit("T_CLOSE_TAG", end)
		t.php = false
		return nil
	}
	if strings.HasPrefix(rest, "/*") {
		n := strings.Index(rest[2:], "*/")
		end := len(t.source)
		if n >= 0 {
			end = start + n + 4
		}
		name := "T_COMMENT"
		if len(rest) > 3 && strings.HasPrefix(rest, "/**") && whitespace(rest[3]) {
			name = "T_DOC_COMMENT"
		}
		t.emit(name, end)
		return nil
	}
	if strings.HasPrefix(rest, "//") || c == '#' && !strings.HasPrefix(rest, "#[") {
		end := start
		for end < len(t.source) && t.source[end] != '\n' && t.source[end] != '\r' && !strings.HasPrefix(t.source[end:], "?>") {
			end++
		}
		t.emit("T_COMMENT", end)
		return nil
	}
	if strings.HasPrefix(rest, "<<<") {
		if m := tokenHeredoc.FindStringSubmatch(rest); m != nil {
			return t.heredoc(m)
		}
	}
	if c == '\'' || c == '"' || c == '`' || (c == 'b' || c == 'B') && len(rest) > 1 && (rest[1] == '\'' || rest[1] == '"') {
		return t.quoted()
	}
	if c == '$' && len(rest) > 1 && nameStart(rest[1]) {
		end := start + 2
		for end < len(t.source) && namePart(t.source[end]) {
			end++
		}
		t.emit("T_VARIABLE", end)
		return nil
	}
	if nameStart(c) || c == '\\' && len(rest) > 1 && nameStart(rest[1]) {
		end := start
		if c == '\\' {
			end++
		}
		for end < len(t.source) && namePart(t.source[end]) {
			end++
		}
		for end+1 < len(t.source) && t.source[end] == '\\' && nameStart(t.source[end+1]) {
			end += 2
			for end < len(t.source) && namePart(t.source[end]) {
				end++
			}
		}
		text := t.source[start:end]
		low := strings.ToLower(text)
		name := "T_STRING"
		switch {
		case strings.HasPrefix(text, "\\"):
			name = "T_NAME_FULLY_QUALIFIED"
		case strings.HasPrefix(low, "namespace\\"):
			name = "T_NAME_RELATIVE"
		case strings.Contains(text, "\\"):
			name = "T_NAME_QUALIFIED"
		default:
			prev := t.previous()
			if prev != "T_OBJECT_OPERATOR" && prev != "T_NULLSAFE_OBJECT_OPERATOR" && !(t.parse && (prev == "T_FUNCTION" || prev == "T_CONST" || prev == "T_DOUBLE_COLON" && low != "class")) {
				if n := phpKeywords[low]; n != "" && phpTokenConstants(t.version)[n] != 0 {
					name = n
				}
			}
		}
		if name == "T_YIELD" {
			i := end
			for i < len(t.source) && whitespace(t.source[i]) {
				i++
			}
			if i+4 <= len(t.source) && strings.EqualFold(t.source[i:i+4], "from") && (i+4 == len(t.source) || !namePart(t.source[i+4])) {
				name = "T_YIELD_FROM"
				end = i + 4
			}
		}
		if t.version != PHP83 && (name == "T_PUBLIC" || name == "T_PROTECTED" || name == "T_PRIVATE") && strings.HasPrefix(t.source[end:], "(set)") {
			name += "_SET"
			end += 5
		}
		t.emit(name, end)
		return nil
	}
	if digit(c) || c == '.' && len(rest) > 1 && digit(rest[1]) {
		s := tokenNumber.FindString(rest)
		if s != "" {
			name := "T_LNUMBER"
			clean := strings.ReplaceAll(s, "_", "")
			hex := strings.HasPrefix(strings.ToLower(clean), "0x")
			if !hex && strings.ContainsAny(clean, ".eE") {
				name = "T_DNUMBER"
			} else if _, err := strconv.ParseUint(clean, 0, 63); err != nil {
				if _, err := strconv.ParseUint(clean, 10, 63); err != nil {
					name = "T_DNUMBER"
				}
			}
			t.emit(name, start+len(s))
			return nil
		}
	}
	if c == '(' {
		if n := strings.IndexByte(rest, ')'); n > 0 && n < 20 {
			word := strings.ToLower(strings.Trim(rest[1:n], " \t\r\n"))
			casts := map[string]string{"int": "INT", "integer": "INT", "bool": "BOOL", "boolean": "BOOL", "float": "DOUBLE", "double": "DOUBLE", "real": "DOUBLE", "string": "STRING", "binary": "STRING", "array": "ARRAY", "object": "OBJECT", "unset": "UNSET", "void": "VOID"}
			if name := casts[word]; name != "" && phpTokenConstants(t.version)["T_"+name+"_CAST"] != 0 {
				t.emit("T_"+name+"_CAST", start+n+1)
				return nil
			}
		}
	}
	if c == '&' && !strings.HasPrefix(rest, "&&") && !strings.HasPrefix(rest, "&=") {
		i := skipTokenTrivia(t.source, start+1)
		name := "T_AMPERSAND_NOT_FOLLOWED_BY_VAR_OR_VARARG"
		if i < len(t.source) && (t.source[i] == '$' || strings.HasPrefix(t.source[i:], "...")) {
			name = "T_AMPERSAND_FOLLOWED_BY_VAR_OR_VARARG"
		}
		t.emit(name, start+1)
		return nil
	}
	for _, op := range phpOperators {
		if strings.HasPrefix(rest, op[0]) && phpTokenConstants(t.version)[op[1]] != 0 {
			t.emit(op[1], start+len(op[0]))
			return nil
		}
	}
	name := ""
	if c < 32 && !whitespace(c) {
		name = "T_BAD_CHARACTER"
	}
	if c == '\\' {
		name = "T_NS_SEPARATOR"
	}
	t.emit(name, start+1)
	return nil
}

func skipTokenTrivia(s string, i int) int {
	for i < len(s) {
		if whitespace(s[i]) {
			i++
			continue
		}
		if strings.HasPrefix(s[i:], "/*") {
			n := strings.Index(s[i+2:], "*/")
			if n < 0 {
				return len(s)
			}
			i += n + 4
			continue
		}
		if strings.HasPrefix(s[i:], "//") || s[i] == '#' && !strings.HasPrefix(s[i:], "#[") {
			for i < len(s) && s[i] != '\n' && s[i] != '\r' {
				i++
			}
			continue
		}
		break
	}
	return i
}
func (t *phpTokenizer) quoted() error {
	start := t.pos
	quotePos := start
	if t.source[start] == 'b' || t.source[start] == 'B' {
		quotePos++
	}
	quote := t.source[quotePos]
	end := quotePos + 1
	interpolation := false
	for end < len(t.source) {
		c := t.source[end]
		if c == '\\' {
			end += 2
			continue
		}
		if c == quote {
			break
		}
		if quote != '\'' && (c == '$' && end+1 < len(t.source) && (nameStart(t.source[end+1]) || t.source[end+1] == '{') || c == '{' && end+1 < len(t.source) && t.source[end+1] == '$') {
			interpolation = true
			if c == '{' || c == '$' && t.source[end+1] == '{' {
				if c == '$' {
					end++
				}
				end++
				depth := 1
				for end < len(t.source) && depth > 0 {
					ch := t.source[end]
					if ch == '\'' || ch == '"' || ch == '`' {
						q := ch
						end++
						for end < len(t.source) {
							if t.source[end] == '\\' {
								end += 2
								continue
							}
							if t.source[end] == q {
								end++
								break
							}
							end++
						}
						continue
					}
					if ch == '{' {
						depth++
					}
					if ch == '}' {
						depth--
					}
					end++
				}
				continue
			}
		}
		end++
	}
	if end > len(t.source) {
		end = len(t.source)
	}
	if quote == '\'' || quote == '"' && !interpolation && end < len(t.source) {
		if end < len(t.source) {
			end++
		}
		t.emit("T_CONSTANT_ENCAPSED_STRING", end)
		return nil
	}
	t.emit("", quotePos+1)
	if err := t.encapsed(end); err != nil {
		return err
	}
	if t.pos < len(t.source) {
		t.emit("", t.pos+1)
	}
	return nil
}
func (t *phpTokenizer) encapsed(end int) error {
	for t.pos < end {
		start := t.pos
		c := t.source[start]
		if c == '{' && start+1 < end && t.source[start+1] == '$' {
			t.emit("T_CURLY_OPEN", start+1)
			depth := 1
			for t.pos < end && depth > 0 {
				if t.source[t.pos] == '}' {
					depth--
					t.emit("", t.pos+1)
				} else if t.source[t.pos] == '{' {
					depth++
					t.emit("", t.pos+1)
				} else if err := t.one(); err != nil {
					return err
				}
			}
			continue
		}
		if c == '$' && start+1 < end && t.source[start+1] == '{' {
			t.emit("T_DOLLAR_OPEN_CURLY_BRACES", start+2)
			if t.pos < end && nameStart(t.source[t.pos]) {
				j := t.pos + 1
				for j < end && namePart(t.source[j]) {
					j++
				}
				t.emit("T_STRING_VARNAME", j)
			}
			for t.pos < end && t.source[t.pos] != '}' {
				if err := t.one(); err != nil {
					return err
				}
			}
			if t.pos < end {
				t.emit("", t.pos+1)
			}
			continue
		}
		if c == '$' && start+1 < end && nameStart(t.source[start+1]) {
			if err := t.one(); err != nil {
				return err
			}
			if t.pos < end && t.source[t.pos] == '[' {
				t.emit("", t.pos+1)
				for t.pos < end && t.source[t.pos] != ']' {
					before := len(t.tokens)
					if err := t.one(); err != nil {
						return err
					}
					if len(t.tokens) > before && t.tokens[len(t.tokens)-1].name == "T_LNUMBER" {
						t.tokens[len(t.tokens)-1].name = "T_NUM_STRING"
					}
				}
				if t.pos < end {
					t.emit("", t.pos+1)
				}
			}
			if t.pos+2 < end && strings.HasPrefix(t.source[t.pos:], "->") && nameStart(t.source[t.pos+2]) {
				t.emit("T_OBJECT_OPERATOR", t.pos+2)
				if err := t.one(); err != nil {
					return err
				}
			}
			continue
		}
		j := start + 1
		if c == '\\' && j < end {
			j++
		}
		for j < end {
			if t.source[j] == '\\' {
				j += 2
				continue
			}
			if t.source[j] == '$' && j+1 < end && (nameStart(t.source[j+1]) || t.source[j+1] == '{') || t.source[j] == '{' && j+1 < end && t.source[j+1] == '$' {
				break
			}
			j++
		}
		if j > end {
			j = end
		}
		t.emit("T_ENCAPSED_AND_WHITESPACE", j)
	}
	return nil
}
func (t *phpTokenizer) heredoc(m []string) error {
	t.emit("T_START_HEREDOC", t.pos+len(m[0]))
	bodyStart := t.pos
	end := len(t.source)
	labelStart := end
	for line := bodyStart; line < len(t.source); {
		j := line
		for j < len(t.source) && (t.source[j] == ' ' || t.source[j] == '\t') {
			j++
		}
		if strings.HasPrefix(t.source[j:], m[2]) && (j+len(m[2]) == len(t.source) || !namePart(t.source[j+len(m[2])])) {
			end = line
			labelStart = j
			break
		}
		n := strings.IndexByte(t.source[line:], '\n')
		if n < 0 {
			break
		}
		line += n + 1
	}
	if m[1] == "'" {
		if end > t.pos {
			t.emit("T_ENCAPSED_AND_WHITESPACE", end)
		}
	} else if err := t.encapsed(end); err != nil {
		return err
	}
	if labelStart < len(t.source) {
		t.emit("T_END_HEREDOC", labelStart+len(m[2]))
	}
	return nil
}
func (e *evaluator) callTokenizer(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	param := "code"
	kind := "string"
	max := 2
	if name == "token_name" {
		param = "id"
		kind = "int"
		max = 1
	}
	b, p, _, err := bindCallArguments(args, []parameter{{name: param}, {name: "flags"}}, node.line)
	if err != nil {
		return value{}, err
	}
	if !p[0] || len(args) > max {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "() expects 1 to " + strconv.Itoa(max) + " arguments"}
	}
	v, _, ok := e.coerceType(b[0], kind, false, e.strictTypes)
	if !ok {
		return value{}, builtinTypeError(node.line, name, 1, param, kind, b[0])
	}
	constants := phpTokenConstants(e.version)
	if name == "token_name" {
		text := "UNKNOWN"
		for key, n := range constants {
			if int64(n) == v.i && key != "TOKEN_PARSE" && key != "T_PAAMAYIM_NEKUDOTAYIM" {
				text = key
				break
			}
		}
		return value{kind: stringKind, s: text}, nil
	}
	flags := int64(0)
	if p[1] {
		f, _, ok := e.coerceType(b[1], "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, 2, "flags", "int", b[1])
		}
		flags = f.i
	}
	if flags&1 != 0 {
		tokens, err := lex(v.s)
		if err == nil {
			_, err = parseWithVersion(tokens, e.version)
		}
		if err != nil {
			if s, ok := err.(*syntaxError); ok {
				return value{}, &runtimeError{node.line, "ParseError", s.message}
			}
			return value{}, err
		}
	}
	t := phpTokenizer{source: v.s, line: 1, version: e.version, parse: flags&1 != 0}
	if err := t.scan(); err != nil {
		return value{}, err
	}
	out := arrayValue()
	for _, token := range t.tokens {
		item := value{kind: stringKind, s: token.text}
		if token.name != "" {
			id, ok := constants[token.name]
			if !ok {
				return value{}, &runtimeError{node.line, "Error", "gophp: unsupported tokenizer token " + token.name}
			}
			item = arrayValue()
			for i, v := range []value{integer(int64(id)), {kind: stringKind, s: token.text}, integer(int64(token.line))} {
				item.a = item.a.with(arrayKey{index: int64(i)}, v)
			}
		}
		out.a = out.a.with(arrayKey{index: out.a.next}, item)
	}
	return out, nil
}
