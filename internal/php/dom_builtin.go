package php

import (
	"html"
	"regexp"
	"strings"
)

// domNode is the small DOM tree needed by Termwind's HTML renderer. It is not
// an HTML5 parser: unsupported DOM APIs remain undefined methods.
type domNode struct {
	name, text string
	attrs      map[string]string
	attrOrder  []string
	parent     *domNode
	children   []*domNode
	object     *phpObject
	owner      *phpObject
}

var domTokens = regexp.MustCompile(`(?s)<!--.*?-->|<![^>]*>|<[^>]*>|[^<]+`)
var domAttrs = regexp.MustCompile(`([A-Za-z_:][A-Za-z0-9_:.-]*)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+))`)

func domValue(object *phpObject) value { return value{kind: objectKind, o: object} }

func domProperty(object *phpObject, name string, v value) {
	object.properties = object.properties.with(arrayKey{isString: true, text: name}, v)
}

func (e *evaluator) newDOMNode(name, content string, owner *phpObject) *domNode {
	className := "domelement"
	switch name {
	case "#text":
		className = "domtext"
	case "#comment":
		className = "domcomment"
	}
	o := &phpObject{class: e.classes[className], properties: &phpArray{}}
	n := &domNode{name: name, text: content, object: o, owner: owner}
	o.dom = n
	return n
}

func domText(n *domNode) string {
	if n.name == "#text" || n.name == "#comment" {
		return n.text
	}
	var b strings.Builder
	for _, child := range n.children {
		if child.name != "#comment" {
			b.WriteString(domText(child))
		}
	}
	return b.String()
}

func (e *evaluator) syncDOMNode(n *domNode) {
	domProperty(n.object, "nodeName", value{kind: stringKind, s: n.name})
	domProperty(n.object, "nodeValue", value{kind: stringKind, s: domText(n)})
	domProperty(n.object, "previousSibling", value{kind: nullKind})
	domProperty(n.object, "nextSibling", value{kind: nullKind})
	children := arrayValue()
	for i, child := range n.children {
		children.a = children.a.with(arrayKey{index: int64(i)}, domValue(child.object))
		e.syncDOMNode(child)
		if i > 0 {
			domProperty(child.object, "previousSibling", domValue(n.children[i-1].object))
		}
		if i+1 < len(n.children) {
			domProperty(child.object, "nextSibling", domValue(n.children[i+1].object))
		}
	}
	domProperty(n.object, "childNodes", children)
	if n.owner != nil {
		domProperty(n.object, "ownerDocument", domValue(n.owner))
	}
}

func (e *evaluator) loadDOMHTML(document *phpObject, source string) bool {
	root := document.dom
	root.children = nil
	stack := []*domNode{root}
	for _, token := range domTokens.FindAllString(source, -1) {
		if strings.HasPrefix(token, "<!--") {
			child := e.newDOMNode("#comment", strings.TrimSuffix(strings.TrimPrefix(token, "<!--"), "-->"), document)
			parent := stack[len(stack)-1]
			child.parent = parent
			parent.children = append(parent.children, child)
			continue
		}
		if strings.HasPrefix(token, "<!") || strings.HasPrefix(token, "<?") {
			continue
		}
		if strings.HasPrefix(token, "</") {
			name := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(token, "</"), ">")))
			for i := len(stack) - 1; i > 0; i-- {
				if stack[i].name == name {
					stack = stack[:i]
					break
				}
			}
			continue
		}
		if strings.HasPrefix(token, "<") {
			body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(token, "<"), ">"))
			selfClosing := strings.HasSuffix(body, "/")
			body = strings.TrimSpace(strings.TrimSuffix(body, "/"))
			parts := strings.Fields(body)
			if len(parts) == 0 {
				continue
			}
			name := strings.ToLower(parts[0])
			child := e.newDOMNode(name, "", document)
			child.attrs = make(map[string]string)
			for _, match := range domAttrs.FindAllStringSubmatch(body[len(parts[0]):], -1) {
				v := match[2]
				if v == "" {
					v = match[3]
				}
				if v == "" {
					v = match[4]
				}
				key := strings.ToLower(match[1])
				if _, exists := child.attrs[key]; !exists { child.attrOrder = append(child.attrOrder, key) }
				child.attrs[key] = html.UnescapeString(v)
			}
			parent := stack[len(stack)-1]
			child.parent = parent
			parent.children = append(parent.children, child)
			if !selfClosing && !domVoidElement(name) {
				stack = append(stack, child)
			}
			continue
		}
		if strings.TrimSpace(token) == "" {
			continue // LIBXML_NOBLANKS in the Termwind path
		}
		parent := stack[len(stack)-1]
		child := e.newDOMNode("#text", html.UnescapeString(token), document)
		child.parent = parent
		parent.children = append(parent.children, child)
	}
	e.syncDOMNode(root)
	return true
}

func domVoidElement(name string) bool {
	switch name {
	case "br", "hr", "img", "input", "meta", "link", "area", "base", "source", "wbr":
		return true
	}
	return false
}

func domFind(n *domNode, name string, out *[]*domNode) {
	for _, child := range n.children {
		if name == "*" || child.name == strings.ToLower(name) {
			*out = append(*out, child)
		}
		domFind(child, name, out)
	}
}

func domXML(n *domNode) string {
	switch n.name {
	case "#text":
		return html.EscapeString(n.text)
	case "#comment":
		return "<!--" + n.text + "-->"
	}
	var b strings.Builder
	if n.name != "#document" {
		b.WriteString("<" + n.name)
		for _, key := range n.attrOrder {
			val := n.attrs[key]
			b.WriteString(" " + key + "=\"" + html.EscapeString(val) + "\"")
		}
		if len(n.children) == 0 {
			return b.String() + "/>"
		}
		b.WriteString(">")
	}
	for _, child := range n.children {
		b.WriteString(domXML(child))
	}
	if n.name != "#document" {
		b.WriteString("</" + n.name + ">")
	}
	return b.String()
}

func (e *evaluator) callDOMMethod(object value, node methodCall) (value, bool, error) {
	name := strings.ToLower(node.name)
	class := strings.ToLower(object.o.className())
	if !((class == "domdocument" && (name == "loadhtml" || name == "getelementsbytagname" || name == "savexml")) ||
		(class == "domelement" && name == "getattribute")) {
		return value{}, false, nil
	}
	args, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, true, err
	}
	if len(args) < 1 {
		return value{}, true, &runtimeError{node.line, "ArgumentCountError", object.o.className() + "::" + node.name + "() expects at least 1 argument"}
	}
	first := args[0].value
	if class == "domdocument" && name == "loadhtml" {
		if first.kind != stringKind {
			return value{}, true, builtinTypeError(node.line, "DOMDocument::loadHTML", 1, "source", "string", first)
		}
		return boolean(e.loadDOMHTML(object.o, first.s)), true, nil
	}
	if class == "domdocument" && name == "savexml" {
		if first.kind != objectKind || first.o == nil || first.o.dom == nil {
			return value{}, true, builtinTypeError(node.line, "DOMDocument::saveXML", 1, "node", "?DOMNode", first)
		}
		return value{kind: stringKind, s: domXML(first.o.dom)}, true, nil
	}
	if first.kind != stringKind {
		return value{}, true, builtinTypeError(node.line, object.o.className()+"::"+node.name, 1, "qualifiedName", "string", first)
	}
	if name == "getattribute" {
		return value{kind: stringKind, s: object.o.dom.attrs[strings.ToLower(first.s)]}, true, nil
	}
	var found []*domNode
	domFind(object.o.dom, first.s, &found)
	list := &phpObject{class: e.classes["domnodelist"], properties: &phpArray{}}
	items := arrayValue()
	for i, n := range found {
		items.a = items.a.with(arrayKey{index: int64(i)}, domValue(n.object))
	}
	list.iteratorArray = items.a
	domProperty(list, "length", integer(int64(len(found))))
	return domValue(list), true, nil
}
