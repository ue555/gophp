package php

import (
	"net/url"
	"strings"
)

func (e *evaluator) initializeRequestArrays(request HTTPRequest) {
	get := arrayValue()
	if uri, err := url.ParseRequestURI(request.URI); err == nil {
		get = parseQueryArray(uri.RawQuery)
	}
	post := arrayValue()
	contentType := ""
	for name, content := range request.Headers {
		if strings.EqualFold(name, "Content-Type") {
			contentType = content
		}
	}
	mediaType, _, _ := strings.Cut(contentType, ";")
	if strings.EqualFold(request.Method, "POST") && strings.EqualFold(strings.TrimSpace(mediaType), "application/x-www-form-urlencoded") {
		post = parseQueryArray(request.Body)
	}
	cookies := arrayValue()
	for name, content := range request.Headers {
		if !strings.EqualFold(name, "Cookie") {
			continue
		}
		for _, part := range strings.Split(content, ";") {
			piece := strings.SplitN(strings.TrimSpace(part), "=", 2)
			if len(piece) == 2 {
				decoded, err := url.QueryUnescape(piece[1])
				if err != nil {
					decoded = piece[1]
				}
				key := arrayKey{isString: true, text: piece[0]}
				if _, exists := cookies.a.get(key); !exists {
					cookies.a = cookies.a.with(key, value{kind: stringKind, s: decoded})
				}
			}
		}
	}
	all := arrayValue()
	for _, part := range []*phpArray{get.a, post.a, cookies.a} {
		for _, entry := range part.entries {
			all.a = all.a.with(entry.key, entry.value)
		}
	}
	e.rootScope["_GET"] = get
	e.rootScope["_POST"] = post
	e.rootScope["_COOKIE"] = cookies
	e.rootScope["_REQUEST"] = all
}
