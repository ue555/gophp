package php

// Attributes are retained with their declaration. Their arguments are only
// evaluated when ReflectionAttribute::getArguments/newInstance is requested.
type attributeDecl struct {
	name string
	args []expr
	line int
}

func (p *parser) attributes() ([]attributeDecl, error) {
	var all []attributeDecl
	for p.accept("#[") {
		for {
			name, err := p.readName()
			if err != nil {
				return nil, err
			}
			resolved, _ := p.resolveName(name.text, "class")
			attr := attributeDecl{name: resolved, line: name.line}
			if p.accept("(") {
				attr.args, err = p.arguments()
				if err != nil {
					return nil, err
				}
			}
			all = append(all, attr)
			if !p.accept(",") {
				break
			}
		}
		if err := p.expect("]"); err != nil {
			return nil, err
		}
	}
	return all, nil
}
