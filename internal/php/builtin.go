package php

import (
	"fmt"
	"strings"
)

func builtinName(name string) bool {
	if streamBuiltin(name) {
		return true
	}
	if strings.EqualFold(name, "stream_select") {
		return true
	}
	if includePathBuiltin(name) {
		return true
	}
	if constantBuiltin(name) {
		return true
	}
	if strings.EqualFold(name, "dirname") || strings.EqualFold(name, "strtr") || strings.EqualFold(name, "strpos") || strings.EqualFold(name, "strrpos") || strings.EqualFold(name, "stripos") || strings.EqualFold(name, "substr") || strings.EqualFold(name, "function_exists") || strings.EqualFold(name, "extension_loaded") || strings.EqualFold(name, "class_exists") || strings.EqualFold(name, "interface_exists") || strings.EqualFold(name, "trait_exists") || strings.EqualFold(name, "enum_exists") || strings.EqualFold(name, "file_exists") || strings.EqualFold(name, "is_dir") || strings.EqualFold(name, "is_file") || strings.EqualFold(name, "is_link") {
		return true
	}
	if strings.EqualFold(name, "file_get_contents") || strings.EqualFold(name, "file_put_contents") {
		return true
	}
	if strings.EqualFold(name, "fileinode") {
		return true
	}
	if strings.EqualFold(name, "class_implements") {
		return true
	}
	if strings.EqualFold(name, "iconv") || strings.EqualFold(name, "iconv_strlen") || strings.EqualFold(name, "iconv_substr") {
		return true
	}
	if strings.EqualFold(name, "strcmp") || strings.EqualFold(name, "strcasecmp") {
		return true
	}
	if strings.EqualFold(name, "strspn") {
		return true
	}
	if strings.EqualFold(name, "unpack") {
		return true
	}
	if strings.EqualFold(name, "basename") {
		return true
	}
	if strings.EqualFold(name, "str_repeat") {
		return true
	}
	if strings.EqualFold(name, "str_pad") {
		return true
	}
	if strings.EqualFold(name, "ksort") {
		return true
	}
	if strings.EqualFold(name, "array_keys") {
		return true
	}
	if strings.EqualFold(name, "array_map") {
		return true
	}
	if strings.EqualFold(name, "array_combine") {
		return true
	}
	if strings.EqualFold(name, "array_diff_key") || strings.EqualFold(name, "array_intersect_key") || strings.EqualFold(name, "array_flip") {
		return true
	}
	if strings.EqualFold(name, "array_replace") {
		return true
	}
	if strings.EqualFold(name, "time") {
		return true
	}
	if strings.EqualFold(name, "microtime") {
		return true
	}
	if strings.EqualFold(name, "strcspn") {
		return true
	}
	if strings.EqualFold(name, "filter_var") {
		return true
	}
	if strings.EqualFold(name, "ord") || strings.EqualFold(name, "chr") {
		return true
	}
	if strings.EqualFold(name, "array_splice") {
		return true
	}
	if strings.EqualFold(name, "base64_encode") || strings.EqualFold(name, "base64_decode") {
		return true
	}
	if strings.EqualFold(name, "class_alias") {
		return true
	}
	if strings.EqualFold(name, "strrev") || strings.EqualFold(name, "strrchr") {
		return true
	}
	if strings.EqualFold(name, "phpversion") {
		return true
	}
	if processBuiltin(name) {
		return true
	}
	if outputBufferBuiltin(name) {
		return true
	}
	if strings.EqualFold(name, "strip_tags") || strings.EqualFold(name, "html_entity_decode") {
		return true
	}
	if strings.EqualFold(name, "usleep") {
		return true
	}
	if strings.EqualFold(name, "mb_strimwidth") {
		return true
	}
	if strings.EqualFold(name, "hexdec") || strings.EqualFold(name, "hex2bin") {
		return true
	}
	if strings.EqualFold(name, "number_format") {
		return true
	}
	if urlCodecBuiltin(name) {
		return true
	}
	if strings.EqualFold(name, "parse_str") {
		return true
	}
	if strings.EqualFold(name, "http_build_query") {
		return true
	}
	if strings.EqualFold(name, "bin2hex") {
		return true
	}
	if strings.EqualFold(name, "request_parse_body") {
		return true
	}
	if strings.EqualFold(name, "get_cfg_var") {
		return true
	}
	if strings.EqualFold(name, "debug_backtrace") {
		return true
	}
	if strings.EqualFold(name, "array_walk") || strings.EqualFold(name, "array_walk_recursive") {
		return true
	}
	if strings.EqualFold(name, "parse_url") {
		return true
	}
	if strings.EqualFold(name, "preg_quote") {
		return true
	}
	if strings.EqualFold(name, "error_reporting") {
		return true
	}
	if strings.EqualFold(name, "set_exception_handler") {
		return true
	}
	if strings.EqualFold(name, "restore_exception_handler") {
		return true
	}
	if strings.EqualFold(name, "register_shutdown_function") {
		return true
	}
	if strings.EqualFold(name, "error_get_last") || strings.EqualFold(name, "error_clear_last") {
		return true
	}
	if strings.EqualFold(name, "extract") {
		return true
	}
	if strings.EqualFold(name, "realpath") {
		return true
	}
	if strings.EqualFold(name, "sprintf") || strings.EqualFold(name, "vsprintf") {
		return true
	}
	if strings.EqualFold(name, "str_contains") || strings.EqualFold(name, "str_starts_with") || strings.EqualFold(name, "str_ends_with") {
		return true
	}
	if strings.EqualFold(name, "explode") {
		return true
	}
	if strings.EqualFold(name, "str_replace") {
		return true
	}
	if strings.EqualFold(name, "json_encode") || strings.EqualFold(name, "json_decode") || strings.EqualFold(name, "json_last_error") || strings.EqualFold(name, "json_last_error_msg") {
		return true
	}
	if strings.EqualFold(name, "preg_match") || strings.EqualFold(name, "preg_match_all") || strings.EqualFold(name, "preg_replace") || strings.EqualFold(name, "preg_split") {
		return true
	}
	if strings.EqualFold(name, "min") || strings.EqualFold(name, "max") {
		return true
	}
	if strings.EqualFold(name, "date_default_timezone_get") || strings.EqualFold(name, "date_default_timezone_set") {
		return true
	}
	if strings.EqualFold(name, "mkdir") {
		return true
	}
	switch strings.ToLower(name) {
	case "spl_autoload_register", "spl_autoload_unregister":
		return true
	case "hash", "sha1", "str_split", "strpbrk", "strval", "str_getcsv", "assert", "rawurldecode", "rawurlencode", "htmlspecialchars", "random_bytes", "gmdate", "date", "headers_sent", "headers_list", "header", "header_remove", "http_response_code", "mb_split", "serialize", "unserialize", "clearstatcache", "compact", "glob", "get_declared_classes", "get_declared_interfaces", "get_class_methods", "get_parent_class", "class_parents", "class_uses", "sort", "asort", "arsort", "usort", "uasort", "uksort", "array_diff", "array_intersect", "array_intersect_assoc", "array_udiff_assoc", "array_pad", "array_search", "array_reduce", "array_any", "array_all", "array_fill_keys", "array_merge_recursive", "array_replace_recursive", "array_change_key_case", "array_sum", "array_column", "range", "ucfirst", "lcfirst", "ucwords", "abs", "ceil", "floor", "round", "version_compare", "error_log", "is_subclass_of", "preg_replace_callback":
		return true
	case "is_callable", "call_user_func", "call_user_func_array", "exit", "die":
		return true
	case "set_error_handler", "restore_error_handler", "trigger_error":
		return true
	case "strlen", "count", "is_countable", "trim", "ltrim", "rtrim", "strtolower", "strtoupper", "strstr", "strtok", "implode", "array_merge", "array_values", "array_reverse", "array_filter", "array_find_key", "array_first", "array_last", "array_unique", "array_shift", "array_unshift", "array_push", "array_pop", "array_slice", "array_key_exists", "in_array", "current", "end", "spl_object_hash", "spl_object_id", "getenv", "fopen", "fwrite", "fflush", "fclose", "flock", "stream_get_contents", "is_resource", "get_resource_type", "stream_isatty", "stream_get_meta_data", "stream_set_chunk_size", "get_class", "property_exists", "method_exists", "is_a", "iterator_to_array", "iterator_count", "func_num_args", "func_get_args", "func_get_arg", "gettype", "get_debug_type", "is_int", "is_float", "is_string", "is_bool", "is_null", "is_array", "is_object", "is_numeric", "is_iterable", "is_scalar":
		return true
	}
	return false
}

func builtinNameForVersion(name string, version Version) bool {
	if (strings.EqualFold(name, "array_any") || strings.EqualFold(name, "array_all") || strings.EqualFold(name, "array_find_key")) && version == PHP83 {
		return false
	}
	if (strings.EqualFold(name, "array_first") || strings.EqualFold(name, "array_last")) && version != PHP85 {
		return false
	}
	return builtinName(name)
}

func builtinRedeclaration(name string, version Version) string {
	prefix := "Cannot redeclare "
	if version != PHP83 {
		prefix += "function "
	}
	return prefix + name + "()"
}

func (e *evaluator) callBuiltin(node functionCall) (result value, err error) {
	if strings.EqualFold(node.name, "glob") {
		return e.callGlob(node)
	}
	if strings.EqualFold(node.name, "array_replace_recursive") {
		return e.callArrayReplaceRecursive(node)
	}
	if strings.EqualFold(node.name, "range") {
		return e.callRange(node)
	}
	if strings.EqualFold(node.name, "array_column") {
		return e.callArrayColumn(node)
	}
	if strings.EqualFold(node.name, "array_sum") {
		return e.callArraySum(node)
	}
	if strings.EqualFold(node.name, "round") {
		return e.callRound(node)
	}
	if strings.EqualFold(node.name, "version_compare") {
		return e.callVersionCompare(node)
	}
	if strings.EqualFold(node.name, "ucwords") {
		return e.callUcwords(node)
	}
	if strings.EqualFold(node.name, "abs") {
		return e.callAbs(node)
	}
	if strings.EqualFold(node.name, "ceil") || strings.EqualFold(node.name, "floor") {
		return e.callCeilFloor(node)
	}
	if strings.EqualFold(node.name, "is_countable") {
		return e.callIsCountable(node)
	}
	if strings.EqualFold(node.name, "get_parent_class") {
		return e.callGetParentClass(node)
	}
	if strings.EqualFold(node.name, "class_parents") {
		return e.callClassParents(node)
	}
	if strings.EqualFold(node.name, "class_uses") {
		return e.callClassUses(node)
	}
	if strings.EqualFold(node.name, "mb_split") {
		return e.callMbSplit(node)
	}
	if strings.EqualFold(node.name, "compact") {
		return e.callCompact(node)
	}
	if strings.EqualFold(node.name, "clearstatcache") {
		_, err := e.evalArguments(node.args)
		if err != nil {
			return value{}, err
		}
		if len(node.args) > 2 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "clearstatcache() expects at most 2 arguments"}
		}
		return value{kind: nullKind}, nil
	}
	if strings.EqualFold(node.name, "array_merge_recursive") {
		return e.callArrayMergeRecursive(node)
	}
	if strings.EqualFold(node.name, "file_put_contents") {
		return e.callFilePutContents(node)
	}
	if strings.EqualFold(node.name, "serialize") || strings.EqualFold(node.name, "unserialize") {
		return e.callSerialization(node)
	}
	if strings.EqualFold(node.name, "strpbrk") {
		return e.callStrpbrk(node)
	}
	if strings.EqualFold(node.name, "get_class_methods") {
		return e.callGetClassMethods(node)
	}
	if strings.EqualFold(node.name, "preg_replace_callback") {
		return e.callPregReplaceCallback(node)
	}
	if strings.EqualFold(node.name, "lcfirst") {
		return e.callLcfirst(node)
	}
	if strings.EqualFold(node.name, "is_subclass_of") {
		return e.callIsSubclassOf(node)
	}
	if strings.EqualFold(node.name, "sha1") || strings.EqualFold(node.name, "str_split") {
		return e.callStringPieces(node)
	}
	if strings.EqualFold(node.name, "strval") {
		return e.callStrval(node)
	}
	if strings.EqualFold(node.name, "str_getcsv") {
		return e.callStrGetCSV(node)
	}
	if strings.EqualFold(node.name, "assert") {
		return e.callAssert(node)
	}
	if strings.EqualFold(node.name, "random_bytes") {
		return e.callRandomBytes(node)
	}
	if strings.EqualFold(node.name, "gmdate") || strings.EqualFold(node.name, "date") {
		return e.callDateFormat(node)
	}
	if strings.EqualFold(node.name, "headers_sent") || strings.EqualFold(node.name, "headers_list") || strings.EqualFold(node.name, "header") || strings.EqualFold(node.name, "header_remove") || strings.EqualFold(node.name, "http_response_code") {
		return e.callHeaderBuiltin(node)
	}
	if strings.EqualFold(node.name, "rawurldecode") {
		return e.callRawURLDecode(node)
	}
	if strings.EqualFold(node.name, "rawurlencode") {
		return e.callRawURLEncode(node)
	}
	if strings.EqualFold(node.name, "htmlspecialchars") {
		return e.callHTMLSpecialChars(node)
	}
	if strings.EqualFold(node.name, "strtok") {
		return e.callStrtok(node)
	}
	if strings.EqualFold(node.name, "strstr") {
		return e.callStrstr(node)
	}
	if strings.EqualFold(node.name, "error_log") {
		return e.callErrorLog(node)
	}
	if strings.EqualFold(node.name, "get_declared_classes") || strings.EqualFold(node.name, "get_declared_interfaces") {
		return e.callGetDeclaredTypes(node)
	}
	if strings.EqualFold(node.name, "sort") {
		return e.callSort(node)
	}
	if strings.EqualFold(node.name, "asort") || strings.EqualFold(node.name, "arsort") {
		return e.callSort(node)
	}
	if strings.EqualFold(node.name, "usort") || strings.EqualFold(node.name, "uasort") || strings.EqualFold(node.name, "uksort") {
		return e.callUsort(node)
	}
	if strings.EqualFold(node.name, "array_diff") {
		return e.callArrayDiff(node)
	}
	if strings.EqualFold(node.name, "array_intersect") {
		return e.callArrayIntersect(node)
	}
	if strings.EqualFold(node.name, "array_intersect_assoc") {
		return e.callArrayIntersectAssoc(node)
	}
	if strings.EqualFold(node.name, "array_udiff_assoc") {
		return e.callArrayUDiffAssoc(node)
	}
	if strings.EqualFold(node.name, "class_alias") {
		return e.callClassAlias(node)
	}
	if strings.EqualFold(node.name, "strrev") {
		return e.callStrrev(node)
	}
	if strings.EqualFold(node.name, "strrchr") {
		return e.callStrrchr(node)
	}
	if strings.EqualFold(node.name, "phpversion") {
		return e.callPHPVersion(node)
	}
	if processBuiltin(node.name) {
		return e.callProcessBuiltin(node)
	}
	if outputBufferBuiltin(node.name) {
		return e.callOutputBufferBuiltin(node)
	}
	if strings.EqualFold(node.name, "strip_tags") || strings.EqualFold(node.name, "html_entity_decode") {
		return e.callHTMLStringBuiltin(node)
	}
	if strings.EqualFold(node.name, "usleep") {
		return e.callUSleep(node)
	}
	if strings.EqualFold(node.name, "mb_strimwidth") {
		return e.callMBStrimwidth(node)
	}
	if strings.EqualFold(node.name, "hexdec") {
		return e.callHexdec(node)
	}
	if strings.EqualFold(node.name, "hex2bin") {
		return e.callHex2bin(node)
	}
	if strings.EqualFold(node.name, "number_format") {
		return e.callNumberFormat(node)
	}
	if urlCodecBuiltin(node.name) {
		return e.callURLCodec(node)
	}
	if strings.EqualFold(node.name, "parse_str") {
		return e.callParseStr(node)
	}
	if strings.EqualFold(node.name, "http_build_query") {
		return e.callHTTPBuildQuery(node)
	}
	if strings.EqualFold(node.name, "bin2hex") {
		return e.callBin2hex(node)
	}
	if strings.EqualFold(node.name, "request_parse_body") {
		return e.callRequestParseBody(node)
	}
	if strings.EqualFold(node.name, "debug_backtrace") {
		return e.callDebugBacktrace(node)
	}
	if strings.EqualFold(node.name, "array_pad") {
		return e.callArrayPad(node)
	}
	if strings.EqualFold(node.name, "array_search") {
		return e.callArraySearch(node)
	}
	if strings.EqualFold(node.name, "array_reduce") {
		return e.callArrayReduce(node)
	}
	if strings.EqualFold(node.name, "array_any") || strings.EqualFold(node.name, "array_all") {
		return e.callArrayPredicate(node)
	}
	if strings.EqualFold(node.name, "array_fill_keys") {
		return e.callArrayFillKeys(node)
	}
	if strings.EqualFold(node.name, "array_intersect_key") {
		return e.callArrayIntersectKey(node)
	}
	if strings.EqualFold(node.name, "array_flip") {
		return e.callArrayFlip(node)
	}
	if strings.EqualFold(node.name, "ucfirst") {
		return e.callUcfirst(node)
	}
	if strings.EqualFold(node.name, "hash") {
		return e.callHash(node)
	}
	if strings.EqualFold(node.name, "exit") || strings.EqualFold(node.name, "die") {
		return e.callExit(node)
	}
	if strings.EqualFold(node.name, "spl_autoload_register") || strings.EqualFold(node.name, "spl_autoload_unregister") {
		return e.callAutoloadBuiltin(node)
	}
	if strings.EqualFold(node.name, "dirname") {
		return e.callDirname(node)
	}
	if strings.EqualFold(node.name, "file_exists") {
		return e.callFileExists(node)
	}
	if strings.EqualFold(node.name, "is_link") {
		return e.callIsLink(node)
	}
	if strings.EqualFold(node.name, "file_get_contents") {
		return e.callFileGetContents(node)
	}
	if strings.EqualFold(node.name, "fileinode") {
		return e.callFileInode(node)
	}
	if strings.EqualFold(node.name, "class_implements") {
		return e.callClassImplements(node)
	}
	if strings.EqualFold(node.name, "iconv") {
		return e.callIconv(node)
	}
	if strings.EqualFold(node.name, "iconv_strlen") || strings.EqualFold(node.name, "iconv_substr") {
		return e.callIconvSlice(node)
	}
	if strings.EqualFold(node.name, "strcmp") || strings.EqualFold(node.name, "strcasecmp") {
		return e.callStringCompare(node)
	}
	if strings.EqualFold(node.name, "strspn") {
		return e.callStrspn(node)
	}
	if strings.EqualFold(node.name, "unpack") {
		return e.callUnpack(node)
	}
	if strings.EqualFold(node.name, "basename") {
		return e.callBasename(node)
	}
	if strings.EqualFold(node.name, "str_repeat") {
		return e.callStrRepeat(node)
	}
	if strings.EqualFold(node.name, "str_pad") {
		return e.callStrPad(node)
	}
	if strings.EqualFold(node.name, "ksort") {
		return e.callKsort(node)
	}
	if strings.EqualFold(node.name, "array_keys") {
		return e.callArrayKeys(node)
	}
	if strings.EqualFold(node.name, "array_map") {
		return e.callArrayMap(node)
	}
	if strings.EqualFold(node.name, "array_combine") {
		return e.callArrayCombine(node)
	}
	if strings.EqualFold(node.name, "array_diff_key") {
		return e.callArrayDiffKey(node)
	}
	if strings.EqualFold(node.name, "array_replace") {
		return e.callArrayReplace(node)
	}
	if strings.EqualFold(node.name, "time") {
		return e.callTime(node)
	}
	if strings.EqualFold(node.name, "microtime") {
		return e.callMicrotime(node)
	}
	if strings.EqualFold(node.name, "strcspn") {
		return e.callStrcspn(node)
	}
	if strings.EqualFold(node.name, "filter_var") {
		return e.callFilterVar(node)
	}
	if strings.EqualFold(node.name, "ord") || strings.EqualFold(node.name, "chr") {
		return e.callOrdChr(node)
	}
	if strings.EqualFold(node.name, "array_splice") {
		return e.callArraySplice(node)
	}
	if strings.EqualFold(node.name, "base64_encode") || strings.EqualFold(node.name, "base64_decode") {
		return e.callBase64(node)
	}
	if strings.EqualFold(node.name, "get_cfg_var") {
		return e.callGetCfgVar(node)
	}
	if strings.EqualFold(node.name, "array_walk") || strings.EqualFold(node.name, "array_walk_recursive") {
		return e.callArrayWalk(node)
	}
	if strings.EqualFold(node.name, "parse_url") {
		return e.callParseURL(node)
	}
	if strings.EqualFold(node.name, "preg_quote") {
		return e.callPregQuote(node)
	}
	if strings.EqualFold(node.name, "error_reporting") {
		return e.callErrorReporting(node)
	}
	if strings.EqualFold(node.name, "set_exception_handler") {
		return e.callSetExceptionHandler(node)
	}
	if strings.EqualFold(node.name, "restore_exception_handler") {
		return e.callRestoreExceptionHandler(node)
	}
	if strings.EqualFold(node.name, "register_shutdown_function") {
		return e.callRegisterShutdownFunction(node)
	}
	if strings.EqualFold(node.name, "error_get_last") {
		if len(node.args) != 0 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "error_get_last() expects exactly 0 arguments"}
		}
		if e.lastError.kind == arrayKind {
			return e.lastError, nil
		}
		return value{kind: nullKind}, nil
	}
	if strings.EqualFold(node.name, "error_clear_last") {
		if len(node.args) != 0 {
			return value{}, &runtimeError{node.line, "ArgumentCountError", "error_clear_last() expects exactly 0 arguments"}
		}
		e.lastError = value{kind: nullKind}
		return value{kind: nullKind}, nil
	}
	if strings.EqualFold(node.name, "extract") {
		return e.callExtract(node)
	}
	if strings.EqualFold(node.name, "realpath") {
		return e.callRealpath(node)
	}
	if strings.EqualFold(node.name, "sprintf") || strings.EqualFold(node.name, "vsprintf") {
		return e.callSprintf(node)
	}
	if strings.EqualFold(node.name, "str_contains") || strings.EqualFold(node.name, "str_starts_with") || strings.EqualFold(node.name, "str_ends_with") {
		return e.callStrMatch(node)
	}
	if strings.EqualFold(node.name, "explode") {
		return e.callExplode(node)
	}
	if strings.EqualFold(node.name, "str_replace") {
		return e.callStrReplace(node)
	}
	if strings.EqualFold(node.name, "json_encode") {
		return e.callJSONEncode(node)
	}
	if strings.EqualFold(node.name, "json_decode") || strings.EqualFold(node.name, "json_last_error") || strings.EqualFold(node.name, "json_last_error_msg") {
		return e.callJSONDecode(node)
	}
	if strings.EqualFold(node.name, "preg_match") {
		return e.callPregMatch(node)
	}
	if strings.EqualFold(node.name, "preg_match_all") {
		return e.callPregMatchAll(node)
	}
	if strings.EqualFold(node.name, "preg_split") {
		return e.callPregSplit(node)
	}
	if strings.EqualFold(node.name, "preg_replace") {
		return e.callPregReplace(node)
	}
	if strings.EqualFold(node.name, "min") || strings.EqualFold(node.name, "max") {
		return e.callMinMax(node)
	}
	if strings.EqualFold(node.name, "date_default_timezone_get") || strings.EqualFold(node.name, "date_default_timezone_set") {
		return e.callDefaultTimezone(node)
	}
	if strings.EqualFold(node.name, "mkdir") {
		return e.callMkdir(node)
	}
	if strings.EqualFold(node.name, "is_dir") {
		return e.callIsDir(node)
	}
	if strings.EqualFold(node.name, "is_file") {
		return e.callIsFile(node)
	}
	if strings.EqualFold(node.name, "getenv") {
		return e.callGetenv(node)
	}
	if streamBuiltin(node.name) {
		return e.callStreamBuiltin(node)
	}
	if strings.EqualFold(node.name, "stream_select") {
		return e.callStreamSelect(node)
	}
	if strings.EqualFold(node.name, "strtr") {
		return e.callStrtr(node)
	}
	if strings.EqualFold(node.name, "strrpos") {
		return e.callStrrpos(node)
	}
	if strings.EqualFold(node.name, "strpos") {
		return e.callStrpos(node)
	}
	if strings.EqualFold(node.name, "stripos") {
		return e.callStripos(node)
	}
	if strings.EqualFold(node.name, "strtolower") || strings.EqualFold(node.name, "strtoupper") {
		return e.callCaseString(node)
	}
	if strings.EqualFold(node.name, "substr") {
		return e.callSubstr(node)
	}
	if strings.EqualFold(node.name, "rtrim") || strings.EqualFold(node.name, "ltrim") || strings.EqualFold(node.name, "trim") {
		return e.callTrim(node)
	}
	if strings.EqualFold(node.name, "implode") {
		return e.callImplode(node)
	}
	if strings.EqualFold(node.name, "get_class") {
		return e.callGetClass(node)
	}
	if strings.EqualFold(node.name, "property_exists") {
		return e.callPropertyExists(node)
	}
	if strings.EqualFold(node.name, "method_exists") {
		return e.callMethodExists(node)
	}
	if strings.EqualFold(node.name, "is_a") {
		return e.callIsA(node)
	}
	if strings.EqualFold(node.name, "iterator_to_array") {
		return e.callIteratorToArray(node)
	}
	if strings.EqualFold(node.name, "iterator_count") {
		return e.callIteratorCount(node)
	}
	if strings.EqualFold(node.name, "array_merge") {
		return e.callArrayMerge(node)
	}
	if strings.EqualFold(node.name, "array_values") {
		return e.callArrayValues(node)
	}
	if strings.EqualFold(node.name, "array_filter") || strings.EqualFold(node.name, "array_unique") {
		return e.callArrayFilterOrUnique(node)
	}
	if strings.EqualFold(node.name, "array_find_key") {
		return e.callArrayFindKey(node)
	}
	if strings.EqualFold(node.name, "array_first") || strings.EqualFold(node.name, "array_last") {
		return e.callArrayEdgeValue(node)
	}
	if strings.EqualFold(node.name, "array_reverse") {
		return e.callArrayReverse(node)
	}
	if strings.EqualFold(node.name, "array_unshift") {
		return e.callArrayUnshift(node)
	}
	if strings.EqualFold(node.name, "array_push") {
		return e.callArrayPush(node)
	}
	if strings.EqualFold(node.name, "array_change_key_case") {
		return e.callArrayChangeKeyCase(node)
	}
	if strings.EqualFold(node.name, "array_shift") {
		return e.callArrayShift(node)
	}
	if strings.EqualFold(node.name, "array_pop") {
		return e.callArrayPop(node)
	}
	if strings.EqualFold(node.name, "in_array") {
		return e.callInArray(node)
	}
	if strings.EqualFold(node.name, "func_num_args") || strings.EqualFold(node.name, "func_get_args") || strings.EqualFold(node.name, "func_get_arg") {
		return e.callFuncArgs(node)
	}
	if strings.EqualFold(node.name, "end") {
		return e.callArrayEnd(node)
	}
	if strings.EqualFold(node.name, "current") {
		return e.callArrayCurrent(node)
	}
	if strings.EqualFold(node.name, "array_slice") {
		return e.callArraySlice(node)
	}
	if strings.EqualFold(node.name, "array_key_exists") {
		return e.callArrayKeyExists(node)
	}
	if strings.EqualFold(node.name, "spl_object_hash") {
		return e.callObjectHash(node)
	}
	if strings.EqualFold(node.name, "spl_object_id") {
		return e.callObjectHash(node)
	}
	if strings.EqualFold(node.name, "function_exists") {
		return e.callFunctionExists(node)
	}
	if strings.EqualFold(node.name, "extension_loaded") {
		return e.callExtensionLoaded(node)
	}
	if strings.EqualFold(node.name, "class_exists") || strings.EqualFold(node.name, "interface_exists") || strings.EqualFold(node.name, "trait_exists") || strings.EqualFold(node.name, "enum_exists") {
		return e.callClassExists(node)
	}
	if strings.EqualFold(node.name, "is_callable") || strings.EqualFold(node.name, "call_user_func") || strings.EqualFold(node.name, "call_user_func_array") {
		return e.callCallableBuiltin(node)
	}
	if includePathBuiltin(node.name) {
		return e.callIncludePathBuiltin(node)
	}
	if constantBuiltin(node.name) {
		return e.callConstantBuiltin(node)
	}
	name := strings.ToLower(node.name)
	if name == "set_error_handler" || name == "restore_error_handler" || name == "trigger_error" {
		return e.callErrorBuiltin(node)
	}
	evaluated, err := e.evalArguments(node.args)
	if err != nil {
		return value{}, err
	}
	params := []parameter{{name: "value"}}
	if name == "strlen" {
		params[0].name = "string"
	}
	if name == "count" {
		params = append(params, parameter{name: "mode"})
	}
	bound, present, _, err := bindCallArguments(evaluated, params, node.line)
	if err != nil {
		return value{}, err
	}
	args := argumentValues(evaluated)
	// Zend compiles single-argument calls to these functions to dedicated
	// opcodes. Arity failures and count($a, $mode) retain an internal call frame.
	framed := len(args) != 1
	if !framed {
		node.line = expressionLine(node.args[0])
	}
	defer func() {
		if err != nil && framed && e.failureFrames == nil {
			e.captureFailure(callFrame{method: name, line: node.line, filename: e.filename, args: args})
		}
	}()
	max := 1
	if name == "count" {
		max = 2
	}
	if len(args) < 1 || len(args) > max || !present[0] {
		limit, count := "exactly", 1
		if name == "count" {
			limit = "at least"
			if len(args) > 2 {
				limit, count = "at most", 2
			}
		}
		plural := ""
		if count != 1 {
			plural = "s"
		}
		return value{}, &runtimeError{node.line, "ArgumentCountError", fmt.Sprintf("%s() expects %s %d argument%s, %d given", name, limit, count, plural, len(args))}
	}
	v := bound[0]
	switch name {
	case "strlen":
		if v.kind == arrayKind || v.kind == objectKind {
			return value{}, builtinTypeError(node.line, name, 1, "string", "string", v)
		}
		if v.kind == nullKind {
			if err := e.diagnostic("Deprecated", node.line, "strlen(): Passing null to parameter #1 ($string) of type string is deprecated"); err != nil {
				return value{}, err
			}
		}
		s, err := e.toText(v, node.line)
		return integer(int64(len(s))), err
	case "count":
		mode := int64(0)
		if present[1] {
			mode, err = e.countMode(bound[1], node.line)
			if err != nil {
				return value{}, err
			}
			if mode != 0 && mode != 1 {
				return value{}, &runtimeError{node.line, "ValueError", "count(): Argument #2 ($mode) must be either COUNT_NORMAL or COUNT_RECURSIVE"}
			}
		}
		if v.kind == objectKind && classImplements(v.o.class, "Countable") {
			if v.o.iteratorArray != nil {
				return integer(int64(len(v.o.iteratorArray.entries))), nil
			}
			method := e.lookupMethod(v.o.class, "count")
			if method == nil {
				return value{}, &runtimeError{node.line, "Error", "Call to undefined method " + v.o.className() + "::count()"}
			}
			return e.invoke(v, method, nil, node.line)
		}
		if v.kind != arrayKind {
			return value{}, builtinTypeError(node.line, name, 1, "value", "Countable|array", v)
		}
		return integer(arrayCount(v.a, mode == 1)), nil
	case "gettype":
		if v.kind == resourceKind {
			return value{kind: stringKind, s: "resource"}, nil
		}
		return value{kind: stringKind, s: [...]string{"NULL", "integer", "double", "string", "boolean", "array", "object"}[v.kind]}, nil
	case "get_debug_type":
		if v.kind == resourceKind {
			if v.stream == nil || v.stream.closed {
				return value{kind: stringKind, s: "resource (closed)"}, nil
			}
			return value{kind: stringKind, s: "resource (stream)"}, nil
		}
		return value{kind: stringKind, s: v.typeName()}, nil
	case "is_numeric":
		if v.kind == intKind || v.kind == floatKind {
			return boolean(true), nil
		}
		if v.kind != stringKind {
			return boolean(false), nil
		}
		_, trailing, ok := v.numeric()
		return boolean(ok && !trailing), nil
	case "is_iterable":
		return boolean(v.kind == arrayKind || v.kind == objectKind && classImplements(v.o.class, "Traversable")), nil
	case "is_scalar":
		return boolean(v.kind == intKind || v.kind == floatKind || v.kind == stringKind || v.kind == boolKind), nil
	case "is_resource":
		return boolean(v.kind == resourceKind && v.stream != nil && !v.stream.closed), nil
	default:
		kind := map[string]valueKind{"is_int": intKind, "is_float": floatKind, "is_string": stringKind, "is_bool": boolKind, "is_null": nullKind, "is_array": arrayKind, "is_object": objectKind}[name]
		return boolean(v.kind == kind), nil
	}
}

func builtinTypeError(line int, name string, number int, parameter, expected string, v value) *runtimeError {
	return &runtimeError{line, "TypeError", fmt.Sprintf("%s(): Argument #%d ($%s) must be of type %s, %s given", name, number, parameter, expected, v.diagnosticName())}
}

func (e *evaluator) countMode(v value, line int) (int64, error) {
	n, trailing, ok := v.numeric()
	if !ok || trailing || n.kind == floatKind && floatOutsideInt(n.f) {
		return 0, builtinTypeError(line, "count", 2, "mode", "int", v)
	}
	if v.kind == nullKind {
		e.diagnostic("Deprecated", line, "count(): Passing null to parameter #2 ($mode) of type int is deprecated")
	}
	if n.kind == intKind {
		return n.i, nil
	}
	i := floatToInt(n.f)
	if float64(i) != n.f {
		source := "float " + floatText(n.f, -1)
		if v.kind == stringKind {
			source = "float-string \"" + v.s + "\""
		}
		e.diagnostic("Deprecated", line, "Implicit conversion from "+source+" to int loses precision")
	}
	return i, nil
}

func arrayCount(a *phpArray, recursive bool) int64 {
	n := int64(len(a.entries))
	if recursive {
		for _, entry := range a.entries {
			if v := dereference(entry.value); v.kind == arrayKind {
				n += arrayCount(v.a, true)
			}
		}
	}
	return n
}

func (e *evaluator) captureFailure(extra ...callFrame) {
	if e.failureFrames != nil {
		return
	}
	e.failureFrames = append(e.snapshotCallFrames(), extra...)
}

func (e *evaluator) snapshotCallFrames() []callFrame {
	var frames []callFrame
	for _, frame := range e.frames {
		frames = append(frames, snapshotCallFrame(frame))
	}
	return frames
}

// A trace records the callee's bound argument slots. Named arguments may
// leave default-filled gaps. A variadic parameter contributes its original
// bound elements, not the array that can later change in the local scope.
func snapshotCallFrame(frame callFrame) callFrame {
	if !frame.traceReady {
		frame.args = append([]value(nil), frame.args...)
		frame.locals = nil
		return frame
	}
	variadic := -1
	if n := len(frame.parameters); n > 0 && frame.parameters[n-1].variadic {
		variadic = n - 1
	}
	last := -1
	for i, present := range frame.present {
		if present && i != variadic {
			last = i
		}
	}
	if variadic >= 0 {
		if packed := frame.boundArgs[variadic]; packed.kind == arrayKind && len(packed.a.entries) > 0 {
			last = variadic - 1
		}
	} else if len(frame.callArgs) > len(frame.parameters) {
		last = len(frame.parameters) - 1
	}
	var args []value
	var keys []arrayKey
	for i := 0; i <= last; i++ {
		v := frame.boundArgs[i]
		if current, ok := frame.locals[frame.parameters[i].name]; ok {
			v = dereference(current)
		}
		args = append(args, v)
		keys = append(keys, arrayKey{index: int64(len(keys))})
	}
	if variadic >= 0 {
		if packed := frame.boundArgs[variadic]; packed.kind == arrayKind {
			for _, entry := range packed.a.entries {
				key := entry.key
				if !key.isString {
					key.index = int64(len(keys))
				}
				args = append(args, dereference(entry.value))
				keys = append(keys, key)
			}
		}
	} else if len(frame.callArgs) > len(frame.parameters) {
		for _, arg := range frame.callArgs[len(frame.parameters):] {
			args = append(args, arg.value)
			keys = append(keys, arrayKey{index: int64(len(keys))})
		}
	}
	frame.args, frame.argKeys = args, keys
	frame.locals, frame.callArgs, frame.boundArgs, frame.present = nil, nil, nil, nil
	return frame
}
