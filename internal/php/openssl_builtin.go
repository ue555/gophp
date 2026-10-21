package php

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
)

func opensslBuiltin(name string) bool {
	switch strings.ToLower(name) {
	case "openssl_cipher_iv_length", "openssl_encrypt", "openssl_decrypt":
		return true
	}
	return false
}
func (e *evaluator) callOpenSSL(node functionCall) (value, error) {
	name := strings.ToLower(node.name)
	encrypt := name == "openssl_encrypt"
	ivOnly := name == "openssl_cipher_iv_length"
	args, err := e.evalArgumentsWithCapture(node.args, func(x expr, n string, named bool, i int) (value, *referencePlace, bool, error) {
		if encrypt && (named && n == "tag" || !named && i == 5) && referenceTargetSyntax(x) && !hasNullsafe(x) {
			p, err := e.referencePlace(x, true)
			if err != nil {
				return value{}, nil, true, err
			}
			return p.read(), &p, true, nil
		}
		return value{}, nil, false, nil
	})
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "data"}, {name: "cipher_algo"}, {name: "passphrase"}, {name: "options"}, {name: "iv"}, {name: "tag"}, {name: "aad"}, {name: "tag_length"}}
	required := 3
	if ivOnly {
		params = params[1:2]
		required = 1
	} else if !encrypt {
		params = params[:7]
	}
	b, p, _, err := bindCallArguments(args, params, node.line)
	if err != nil {
		return value{}, err
	}
	if len(args) < required || len(args) > len(params) {
		return value{}, &runtimeError{node.line, "ArgumentCountError", name + "(): Invalid argument count"}
	}
	for i := range b {
		if !p[i] {
			if i < required {
				return value{}, &runtimeError{node.line, "ArgumentCountError", name + "(): Missing argument"}
			}
			continue
		}
		if !ivOnly && i == 5 && encrypt {
			continue
		}
		typ := "string"
		if !ivOnly && (i == 3 || i == 7) {
			typ = "int"
		}
		v, _, ok := e.coerceType(b[i], typ, false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(node.line, name, i+1, params[i].name, typ, b[i])
		}
		b[i] = v
	}
	fail := func(message string) (value, error) {
		if message != "" {
			if err := e.warn(node.line, name+"(): "+message); err != nil {
				return value{}, err
			}
		}
		return boolean(false), nil
	}
	algo := b[0].s
	if !ivOnly {
		algo = b[1].s
	}
	algo = strings.ToLower(algo)
	sizes := map[string]int{"aes-128-cbc": 16, "aes-192-cbc": 24, "aes-256-cbc": 32, "aes-128-gcm": 16, "aes-192-gcm": 24, "aes-256-gcm": 32}
	size, ok := sizes[algo]
	if !ok {
		return fail("Unknown cipher algorithm")
	}
	gcm := strings.HasSuffix(algo, "-gcm")
	ivlen := 16
	if gcm {
		ivlen = 12
	}
	if ivOnly {
		return integer(int64(ivlen)), nil
	}
	options := b[3].i
	key := []byte(b[2].s)
	if len(key) < size {
		if options&4 != 0 {
			return fail("Key length cannot be set for the cipher algorithm")
		}
		key = append(key, make([]byte, size-len(key))...)
	}
	key = key[:size]
	iv := []byte(b[4].s)
	if !gcm && len(iv) != ivlen {
		message := ""
		if len(iv) == 0 {
			if encrypt {
				message = "Using an empty Initialization Vector (iv) is potentially insecure and not recommended"
			}
		} else if len(iv) < ivlen {
			message = fmt.Sprintf("IV passed is only %d bytes long, cipher expects an IV of precisely %d bytes, padding with \\0", len(iv), ivlen)
		} else {
			message = fmt.Sprintf("IV passed is %d bytes long which is longer than the %d expected by selected cipher, truncating", len(iv), ivlen)
		}
		if message != "" {
			if err := e.warn(node.line, name+"(): "+message); err != nil {
				return value{}, err
			}
		}
		iv = append(iv, make([]byte, ivlen)...)[:ivlen]
	}
	block, _ := aes.NewCipher(key)
	data := []byte(b[0].s)
	if !encrypt && options&1 == 0 {
		data, err = base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return fail("")
		}
	}
	var out []byte
	var tag value
	if gcm {
		if len(iv) == 0 {
			return fail("Setting of IV length for AEAD mode failed")
		}
		taglen := 16
		if encrypt && p[7] {
			taglen = int(b[7].i)
		} else if !encrypt {
			taglen = len(b[5].s)
		}
		if taglen < 12 || taglen > 16 {
			return fail("Setting tag length for AEAD cipher failed")
		}
		var aead cipher.AEAD
		if len(iv) == 12 {
			aead, err = cipher.NewGCMWithTagSize(block, taglen)
		} else if taglen == 16 {
			aead, err = cipher.NewGCMWithNonceSize(block, len(iv))
		} else {
			return value{}, &runtimeError{node.line, "Error", "gophp: GCM custom nonce and truncated tag together are unsupported"}
		}
		if err != nil {
			return fail(err.Error())
		}
		if encrypt {
			sealed := aead.Seal(nil, iv, data, []byte(b[6].s))
			out = sealed[:len(sealed)-taglen]
			tag = value{kind: stringKind, s: string(sealed[len(sealed)-taglen:])}
		} else {
			out, err = aead.Open(nil, iv, append(data, []byte(b[5].s)...), []byte(b[6].s))
			if err != nil {
				return fail("")
			}
		}
	} else if encrypt {
		if options&2 == 0 {
			n := 16 - len(data)%16
			data = append(data, bytes.Repeat([]byte{byte(n)}, n)...)
		} else if len(data)%16 != 0 {
			return fail("")
		}
		out = make([]byte, len(data))
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(out, data)
	} else {
		if len(data)%16 != 0 || len(data) == 0 {
			return fail("")
		}
		out = make([]byte, len(data))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, data)
		if options&2 == 0 {
			n := int(out[len(out)-1])
			if n < 1 || n > 16 || !bytes.Equal(out[len(out)-n:], bytes.Repeat([]byte{byte(n)}, n)) {
				return fail("")
			}
			out = out[:len(out)-n]
		}
	}
	if encrypt && p[5] {
		a := argumentForParameter(args, "tag", 5)
		if a == nil || a.place == nil && a.ref == nil {
			return value{}, &runtimeError{node.line, "Error", name + "(): Argument #6 ($tag) cannot be passed by reference"}
		}
		if a.place != nil {
			err = a.place.assign(tag)
		} else {
			err = e.writeReference(a.ref, tag, node.line)
		}
		if err != nil {
			return value{}, err
		}
	}
	result := string(out)
	if encrypt && options&1 == 0 {
		result = base64.StdEncoding.EncodeToString(out)
	}
	return value{kind: stringKind, s: result}, nil
}
