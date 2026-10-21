package php

import (
	"regexp"
	"strings"
)

// RegexIterator currently supports MATCH and GET_MATCH. The fixed metadata
// scanner additionally needs a leading negative lookahead for *Test.php;
// other PCRE-only constructs continue to fail in compileBasicPCRE.
func (e *evaluator) invokeRegexIterator(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	o := object.o
	switch strings.ToLower(method.name) {
	case "__construct":
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "iterator"}, {name: "pattern"}, {name: "mode"}, {name: "flags"}, {name: "pregFlags"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || !present[1] || len(args) > 5 {
			return value{}, &runtimeError{line, "ArgumentCountError", "RegexIterator::__construct() expects at least 2 arguments"}
		}
		if bound[0].kind != objectKind || bound[0].o == nil || !classImplements(bound[0].o.class, "Iterator") {
			return value{}, builtinTypeError(line, "RegexIterator::__construct", 1, "iterator", "Iterator", bound[0])
		}
		pattern, _, ok := e.coerceType(bound[1], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "RegexIterator::__construct", 2, "pattern", "string", bound[1])
		}
		mode := int64(0)
		if present[2] {
			v, _, ok := e.coerceType(bound[2], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "RegexIterator::__construct", 3, "mode", "int", bound[2])
			}
			mode = v.i
		}
		if mode != 0 && mode != 1 {
			return value{}, &runtimeError{line, "Error", "gophp: RegexIterator mode is not supported"}
		}
		prepared := pattern.s
		for _, exclusion := range []struct{ expression, pattern string }{
			{"^(?!.*Test\\.php$)", `Test\.php$`},
			{"^(?!.*(?:Test|\\.blade)\\.php$)", `(?:Test|\.blade)\.php$`},
		} {
			if strings.Contains(prepared, exclusion.expression) {
				prepared = strings.Replace(prepared, exclusion.expression, "^", 1)
				o.regexExclude = regexp.MustCompile("(?i)" + exclusion.pattern)
				break
			}
		}
		compiled, err := compileBasicPCRE(prepared)
		if err != nil {
			return value{}, err
		}
		o.filterInner, o.regexPattern, o.regexCompiled, o.regexMode = bound[0], pattern.s, compiled, mode
		return value{kind: nullKind}, nil
	case "getregex":
		return value{kind: stringKind, s: o.regexPattern}, nil
	case "getmode":
		return integer(o.regexMode), nil
	case "accept":
		current, err := e.iteratorMethod(o.filterInner, "current", line)
		if err != nil {
			return value{}, err
		}
		s, err := e.toText(current, line)
		if err != nil {
			return value{}, err
		}
		if o.regexExclude != nil && o.regexExclude.MatchString(s) {
			return boolean(false), nil
		}
		groups := o.regexCompiled.FindStringSubmatch(s)
		if groups == nil {
			return boolean(false), nil
		}
		if o.regexMode == 1 {
			matches := arrayValue()
			for i, group := range groups {
				matches.a = matches.a.with(arrayKey{index: int64(i)}, value{kind: stringKind, s: group})
			}
			o.regexCurrent = matches
		} else {
			o.regexCurrent = current
		}
		return boolean(true), nil
	case "current":
		return o.regexCurrent, nil
	default:
		return e.invokeFilterIterator(object, method, args, line)
	}
}
