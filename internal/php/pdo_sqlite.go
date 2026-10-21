package php

/*
#cgo LDFLAGS: -lsqlite3
#include <sqlite3.h>
#include <stdlib.h>
static int gophp_bind_text(sqlite3_stmt *stmt,int index,const char *text,int length) {return sqlite3_bind_text(stmt,index,text,length,SQLITE_TRANSIENT);}
static int gophp_bind_blob(sqlite3_stmt *stmt,int index,const char *data,int length) {return sqlite3_bind_blob(stmt,index,data,length,SQLITE_TRANSIENT);}
*/
import "C"

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"
)

type sqliteConnection struct {
	handle           *C.sqlite3
	lastID           int64
	statements       []*sqliteStatement
	queryStmt        *sqliteStatement
	errorState       string
	errorInfo        *pdoFailure
	errorMode        int64
	caseMode         int64
	nullMode         int64
	stringify        bool
	defaultFetchMode int64
}
type sqliteStatement struct {
	owner       *sqliteConnection
	handle      *C.sqlite3_stmt
	bound       map[pdoBoundKey]pdoBoundValue
	boundOrder  []pdoBoundKey
	executed    bool
	finished    bool
	columns     int
	columnNames []string
	affected    int64
	fetchMode   int64
	errorState  string
	errorInfo   *pdoFailure
}

type pdoBoundValue struct {
	data      value
	paramType int64
}

type pdoBoundKey struct {
	index int
	name  string
}

func (stmt *sqliteStatement) rememberBound(key pdoBoundKey, bound pdoBoundValue) {
	if _, exists := stmt.bound[key]; !exists {
		stmt.boundOrder = append(stmt.boundOrder, key)
	}
	stmt.bound[key] = bound
}

// Keep PDO's SQLSTATE and native driver details alongside the ordinary
// runtime error so catch/trace handling remains on the shared Throwable path.
type pdoFailure struct {
	cause  *runtimeError
	state  string
	native int64
	detail string
	impl   bool
}

func (failure *pdoFailure) Error() string { return failure.cause.Error() }
func (failure *pdoFailure) Unwrap() error { return failure.cause }

func pdoErrorCode(state string) value {
	if state == "" {
		return value{kind: nullKind}
	}
	return value{kind: stringKind, s: state}
}

func pdoErrorDetails(state string, failure *pdoFailure) value {
	info := arrayValue()
	info.a = info.a.with(arrayKey{index: 0}, value{kind: stringKind, s: state})
	native, detail := value{kind: nullKind}, value{kind: nullKind}
	if failure != nil {
		native = integer(failure.native)
		if !failure.impl {
			detail = value{kind: stringKind, s: failure.detail}
		}
	}
	info.a = info.a.with(arrayKey{index: 1}, native)
	if failure == nil || !failure.impl {
		info.a = info.a.with(arrayKey{index: 2}, detail)
	}
	return info
}

func pdoIntAttribute(mode value, line int) (int64, error) {
	mode = dereference(mode)
	var selected int64
	switch mode.kind {
	case intKind, boolKind:
		selected = mode.i
	case stringKind:
		parsed, err := strconv.ParseInt(strings.TrimSpace(mode.s), 10, 64)
		if err != nil {
			return 0, &runtimeError{line, "TypeError", "Attribute value must be of type int for selected attribute, string given"}
		}
		selected = parsed
	default:
		return 0, &runtimeError{line, "TypeError", "Attribute value must be of type int for selected attribute, " + mode.typeName() + " given"}
	}
	return selected, nil
}

func (db *sqliteConnection) setAttribute(attr int64, v value, line int, version Version) (bool, error) {
	switch attr {
	case 2:
		seconds, err := pdoIntAttribute(v, line)
		if err != nil {
			return false, err
		}
		C.sqlite3_busy_timeout(db.handle, C.int(seconds*1000))
		return true, nil
	case 3, 8, 11, 19:
		selected, err := pdoIntAttribute(v, line)
		if err != nil {
			return false, err
		}
		switch attr {
		case 3:
			if selected < 0 || selected > 2 {
				return false, &runtimeError{line, "ValueError", "Error mode must be one of the PDO::ERRMODE_* constants"}
			}
			db.errorMode = selected
		case 8:
			if selected < 0 || selected > 2 {
				return false, &runtimeError{line, "ValueError", "Case folding mode must be one of the PDO::CASE_* constants"}
			}
			db.caseMode = selected
		case 11:
			// PHP 8.3/8.4 truncate this attribute to the two-bit PDO field;
			// PHP 8.5 retains the integer assigned to it.
			if version == PHP85 {
				db.nullMode = selected
			} else {
				db.nullMode = selected & 3
			}
		case 19:
			if selected == 0 {
				return false, &runtimeError{line, "ValueError", "Fetch mode must be a bitmask of PDO::FETCH_* constants"}
			}
			db.defaultFetchMode = selected
		}
		return true, nil
	case 17:
		v = dereference(v)
		switch v.kind {
		case boolKind, intKind:
			db.stringify = v.i != 0
			return true, nil
		default:
			return false, &runtimeError{line, "TypeError", "Attribute value must be of type bool for selected attribute, " + v.typeName() + " given"}
		}
	}
	// SQLite has no driver implementation for the other attributes, including
	// ATTR_EMULATE_PREPARES. Do not claim to have applied an ignored setting.
	return false, nil
}

func sqliteError(line int, db *C.sqlite3) error {
	code := C.sqlite3_errcode(db)
	detail := C.GoString(C.sqlite3_errmsg(db))
	state, description := "HY000", "General error"
	switch code {
	case C.SQLITE_NOTFOUND:
		state, description = "42S02", "Base table or view not found"
	case C.SQLITE_INTERRUPT:
		state, description = "01002", "Disconnect error"
	case C.SQLITE_NOLFS:
		state, description = "HYC00", "Optional feature not implemented"
	case C.SQLITE_TOOBIG:
		state, description = "22001", "String data, right truncated"
	case C.SQLITE_CONSTRAINT:
		state, description = "23000", "Integrity constraint violation"
	}
	message := fmt.Sprintf("SQLSTATE[%s]: %s: %d %s", state, description, int(code), detail)
	return &pdoFailure{cause: &runtimeError{line, "PDOException", message}, state: state, native: int64(code), detail: detail}
}

func (e *evaluator) pdoOperationFailure(db *sqliteConnection, method string, err error, line int) (value, error) {
	if _, ok := err.(*pdoFailure); !ok || db.errorMode == 2 {
		return value{}, err
	}
	if db.errorMode == 1 {
		if !strings.Contains(method, "::") {
			method = "PDO::" + method
		}
		if warning := e.warn(line, method+"(): "+err.Error()); warning != nil {
			return value{}, warning
		}
	}
	return boolean(false), nil
}
func (e *evaluator) invokePDO(object value, method *methodDecl, args []callArgument, line int) (value, error) {
	name := strings.ToLower(method.name)
	if method.owner.name == "PDOStatement" {
		return e.invokePDOStatement(object, name, args, line)
	}
	if name == "connect" || name == "__construct" {
		bound, present, _, err := bindCallArguments(args, []parameter{{name: "dsn"}, {name: "username"}, {name: "password"}, {name: "options"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::__construct() expects at least 1 argument"}
		}
		dsn, _, ok := e.coerceType(bound[0], "string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "PDO::__construct", 1, "dsn", "string", bound[0])
		}
		if !strings.HasPrefix(dsn.s, "sqlite:") {
			return value{}, &runtimeError{line, "PDOException", "gophp: only PDO SQLite is implemented"}
		}
		path := strings.TrimPrefix(dsn.s, "sqlite:")
		cpath := C.CString(path)
		defer C.free(unsafe.Pointer(cpath))
		var db *C.sqlite3
		if C.sqlite3_open(cpath, &db) != C.SQLITE_OK {
			defer C.sqlite3_close(db)
			return value{}, sqliteError(line, db)
		}
		connection := &sqliteConnection{handle: db, errorMode: 2, defaultFetchMode: 4}
		C.sqlite3_busy_timeout(db, 60000)
		if present[3] && bound[3].kind == arrayKind && bound[3].a != nil {
			for _, entry := range bound[3].a.entries {
				if entry.key.isString {
					continue
				}
				if _, err := connection.setAttribute(entry.key.index, entry.value, line, e.version); err != nil {
					C.sqlite3_close(db)
					return value{}, err
				}
			}
		}
		e.sqliteConnections = append(e.sqliteConnections, connection)
		if name == "connect" {
			e.initClasses()
			return value{kind: objectKind, o: &phpObject{class: e.classes["pdo"], properties: &phpArray{}, pdo: connection}}, nil
		}
		object.o.pdo = connection
		return value{kind: nullKind}, nil
	}
	if object.kind != objectKind || object.o == nil || object.o.pdo == nil {
		return value{}, &runtimeError{line, "Error", "Call to undefined PDO connection"}
	}
	db := object.o.pdo
	switch name {
	case "errorcode", "errorinfo":
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::" + method.name + "() expects exactly 0 arguments"}
		}
		state, failure := db.errorState, db.errorInfo
		if db.queryStmt != nil {
			state, failure = db.queryStmt.errorState, db.queryStmt.errorInfo
		}
		if name == "errorcode" {
			return pdoErrorCode(state), nil
		}
		return pdoErrorDetails(state, failure), nil
	case "getattribute":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::getAttribute() expects exactly 1 argument"}
		}
		attribute, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "PDO::getAttribute", 1, "attribute", "int", args[0].value)
		}
		db.errorState, db.errorInfo, db.queryStmt = "00000", nil, nil
		if attribute.i == 3 {
			return integer(db.errorMode), nil
		}
		switch attribute.i {
		case 4, 5:
			return value{kind: stringKind, s: C.GoString(C.sqlite3_libversion())}, nil
		case 8:
			return integer(db.caseMode), nil
		case 11:
			return integer(db.nullMode), nil
		case 16:
			return value{kind: stringKind, s: "sqlite"}, nil
		case 17:
			return boolean(db.stringify), nil
		case 19:
			return integer(db.defaultFetchMode), nil
		}
		failure := &pdoFailure{cause: &runtimeError{line, "PDOException", "SQLSTATE[IM001]: Driver does not support this function: driver does not support that attribute"}, state: "IM001", impl: true}
		db.errorState, db.errorInfo = failure.state, failure
		return e.pdoOperationFailure(db, method.name, failure, line)
	case "setattribute":
		if len(args) != 2 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::setAttribute() expects exactly 2 arguments"}
		}
		attribute, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "PDO::setAttribute", 1, "attribute", "int", args[0].value)
		}
		db.errorState, db.errorInfo, db.queryStmt = "00000", nil, nil
		ok, err := db.setAttribute(attribute.i, args[1].value, line, e.version)
		if err != nil {
			return value{}, err
		}
		return boolean(ok), nil
	case "intransaction":
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::inTransaction() expects exactly 0 arguments"}
		}
		return boolean(C.sqlite3_get_autocommit(db.handle) == 0), nil
	case "begintransaction", "commit", "rollback":
		query := "BEGIN"
		if name == "commit" {
			query = "COMMIT"
		}
		if name == "rollback" {
			query = "ROLLBACK"
		}
		if err := db.execute(query, line); err != nil {
			return e.pdoOperationFailure(db, method.name, err, line)
		}
		return boolean(true), nil
	case "lastinsertid":
		return value{kind: stringKind, s: fmt.Sprint(int64(C.sqlite3_last_insert_rowid(db.handle)))}, nil
	case "quote":
		if len(args) < 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::quote() expects at least 1 argument"}
		}
		s, err := e.toText(args[0].value, line)
		if err != nil {
			return value{}, err
		}
		return value{kind: stringKind, s: "'" + strings.ReplaceAll(s, "'", "''") + "'"}, nil
	case "exec":
		if len(args) != 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::exec() expects exactly 1 argument"}
		}
		s, err := e.toText(args[0].value, line)
		if err != nil {
			return value{}, err
		}
		if err := db.execute(s, line); err != nil {
			return e.pdoOperationFailure(db, method.name, err, line)
		}
		return integer(int64(C.sqlite3_changes(db.handle))), nil
	case "prepare", "query":
		if len(args) < 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDO::" + method.name + "() expects at least 1 argument"}
		}
		s, err := e.toText(args[0].value, line)
		if err != nil {
			return value{}, err
		}
		statement, err := db.prepare(s, line)
		if err != nil {
			return e.pdoOperationFailure(db, method.name, err, line)
		}
		if name == "query" {
			if err := e.executeSQLiteStatement(statement, nil, line); err != nil {
				db.queryStmt = statement
				return e.pdoOperationFailure(db, method.name, err, line)
			}
		}
		e.initClasses()
		return value{kind: objectKind, o: &phpObject{class: e.classes["pdostatement"], properties: &phpArray{}, pdoStatement: statement}}, nil
	}
	return value{}, &runtimeError{line, "Error", "Call to undefined method PDO::" + method.name + "()"}
}

func (db *sqliteConnection) execute(sql string, line int) error {
	db.queryStmt = nil
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var message *C.char
	if C.sqlite3_exec(db.handle, csql, nil, nil, &message) != C.SQLITE_OK {
		if message != nil {
			defer C.sqlite3_free(unsafe.Pointer(message))
		}
		failure := sqliteError(line, db.handle).(*pdoFailure)
		db.errorState, db.errorInfo = failure.state, failure
		return failure
	}
	db.errorState, db.errorInfo = "00000", nil
	return nil
}
func (db *sqliteConnection) prepare(sql string, line int) (*sqliteStatement, error) {
	db.queryStmt = nil
	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))
	var handle *C.sqlite3_stmt
	if C.sqlite3_prepare_v2(db.handle, csql, -1, &handle, nil) != C.SQLITE_OK {
		failure := sqliteError(line, db.handle).(*pdoFailure)
		db.errorState, db.errorInfo = failure.state, failure
		return nil, failure
	}
	db.errorState, db.errorInfo = "00000", nil
	statement := &sqliteStatement{owner: db, handle: handle, bound: make(map[pdoBoundKey]pdoBoundValue), columns: int(C.sqlite3_column_count(handle)), fetchMode: db.defaultFetchMode}
	db.statements = append(db.statements, statement)
	return statement, nil
}

// Each interpreter run owns its PDO handles. Finalize open cursors before
// closing their connections so later HTTP requests cannot inherit read locks.
func (e *evaluator) closeSQLiteConnections() {
	for _, db := range e.sqliteConnections {
		for _, stmt := range db.statements {
			if stmt.handle != nil {
				C.sqlite3_finalize(stmt.handle)
				stmt.handle = nil
			}
		}
		if db.handle != nil {
			C.sqlite3_close(db.handle)
			db.handle = nil
		}
	}
}

// PDO SQLite's parameter hook applies the declared PDO type at execute time.
// Keep this separate from execute(array), whose existing binding path has no
// explicit bindValue type declaration.
func (e *evaluator) bindPDOValue(stmt *sqliteStatement, index int, bound *pdoBoundValue, line int) error {
	v := dereference(bound.data)
	var result C.int
	switch bound.paramType & 0xff {
	case 0: // PDO::PARAM_NULL binds null regardless of the supplied PHP value.
		result = C.sqlite3_bind_null(stmt.handle, C.int(index))
	case 1, 5: // PDO::PARAM_INT, PDO::PARAM_BOOL
		if v.kind == nullKind {
			result = C.sqlite3_bind_null(stmt.handle, C.int(index))
		} else {
			converted, err := e.cast(v, "int", line)
			if err != nil {
				return err
			}
			bound.data = converted
			result = C.sqlite3_bind_int64(stmt.handle, C.int(index), C.sqlite3_int64(converted.i))
		}
	case 2, 3: // PDO::PARAM_STR, PDO::PARAM_LOB
		if v.kind == nullKind {
			result = C.sqlite3_bind_null(stmt.handle, C.int(index))
		} else {
			var data string
			if bound.paramType&0xff == 3 && v.kind == resourceKind {
				if v.stream == nil || v.stream.file == nil || v.stream.closed {
					return &runtimeError{line, "PDOException", "SQLSTATE[HY105]: Invalid parameter type: Expected a stream resource"}
				}
				bytes, err := io.ReadAll(v.stream.file)
				if err != nil {
					return &runtimeError{line, "PDOException", "SQLSTATE[HY105]: Invalid parameter type: Expected a stream resource"}
				}
				data = string(bytes)
				bound.data = value{kind: stringKind, s: data}
			} else {
				converted, err := e.toText(v, line)
				if err != nil {
					return err
				}
				data = converted
				bound.data = value{kind: stringKind, s: data}
			}
			cdata := C.CString(data)
			if bound.paramType&0xff == 3 {
				result = C.gophp_bind_blob(stmt.handle, C.int(index), cdata, C.int(len(data)))
			} else {
				result = C.gophp_bind_text(stmt.handle, C.int(index), cdata, C.int(len(data)))
			}
			C.free(unsafe.Pointer(cdata))
		}
	case 4: // PDO::PARAM_STMT is rejected by the SQLite parameter hook.
		return &runtimeError{line, "PDOException", "gophp: PDO::PARAM_STMT is not supported by SQLite"}
	default: // SQLite PDO uses its text path for unrecognized parameter types.
		textBound := pdoBoundValue{data: v, paramType: 2}
		if err := e.bindPDOValue(stmt, index, &textBound, line); err != nil {
			return err
		}
		bound.data = textBound.data
		return nil
	}
	if result != C.SQLITE_OK {
		return sqliteError(line, stmt.owner.handle)
	}
	return nil
}
func (e *evaluator) executeSQLiteStatement(stmt *sqliteStatement, params *phpArray, line int) error {
	C.sqlite3_reset(stmt.handle)
	C.sqlite3_clear_bindings(stmt.handle)
	stmt.errorState, stmt.errorInfo = "00000", nil
	if params != nil {
		// execute(array) replaces prior bindValue bindings and registers every
		// supplied element with PDO::PARAM_STR.
		stmt.bound = make(map[pdoBoundKey]pdoBoundValue)
		stmt.boundOrder = nil
		for _, entry := range params.entries {
			key := pdoBoundKey{index: int(entry.key.index) + 1}
			if entry.key.isString {
				key.name = entry.key.text
				if !strings.HasPrefix(key.name, ":") {
					key.name = ":" + key.name
				}
				key.index = 0
			}
			boundValue := dereference(entry.value)
			if boundValue.kind != nullKind {
				converted, err := e.toText(boundValue, line)
				if err != nil {
					stmt.recordError(err)
					return err
				}
				boundValue = value{kind: stringKind, s: converted}
			}
			bound := pdoBoundValue{data: boundValue, paramType: 2}
			stmt.rememberBound(key, bound)
		}
	}
	for _, key := range stmt.boundOrder {
		index := key.index
		if key.name != "" {
			cname := C.CString(key.name)
			index = int(C.sqlite3_bind_parameter_index(stmt.handle, cname))
			C.free(unsafe.Pointer(cname))
		}
		bound := stmt.bound[key]
		if err := e.bindPDOValue(stmt, index, &bound, line); err != nil {
			stmt.recordError(err)
			return err
		}
		stmt.bound[key] = bound
	}
	stmt.executed = true
	stmt.finished = false
	if stmt.columns > 0 && stmt.columnNames == nil {
		stmt.columnNames = make([]string, stmt.columns)
		for i := range stmt.columnNames {
			name := C.GoString(C.sqlite3_column_name(stmt.handle, C.int(i)))
			if stmt.owner.caseMode == 1 {
				name = strings.ToUpper(name)
			} else if stmt.owner.caseMode == 2 {
				name = strings.ToLower(name)
			}
			stmt.columnNames[i] = name
		}
	}
	if stmt.columns == 0 {
		result := C.sqlite3_step(stmt.handle)
		if result != C.SQLITE_DONE {
			err := sqliteError(line, stmt.owner.handle)
			stmt.recordError(err)
			return err
		}
		stmt.finished = true
		stmt.affected = int64(C.sqlite3_changes(stmt.owner.handle))
	}
	stmt.errorState, stmt.errorInfo = "00000", nil
	return nil
}

func (stmt *sqliteStatement) recordError(err error) {
	if failure, ok := err.(*pdoFailure); ok {
		stmt.errorState, stmt.errorInfo = failure.state, failure
	}
}
func (stmt *sqliteStatement) columnValue(index int) value {
	columnType := C.sqlite3_column_type(stmt.handle, C.int(index))
	if stmt.owner.stringify && (columnType == C.SQLITE_INTEGER || columnType == C.SQLITE_FLOAT) {
		pointer := C.sqlite3_column_text(stmt.handle, C.int(index))
		length := C.sqlite3_column_bytes(stmt.handle, C.int(index))
		return value{kind: stringKind, s: C.GoStringN((*C.char)(unsafe.Pointer(pointer)), length)}
	}
	var result value
	switch columnType {
	case C.SQLITE_INTEGER:
		result = integer(int64(C.sqlite3_column_int64(stmt.handle, C.int(index))))
	case C.SQLITE_FLOAT:
		result = floating(float64(C.sqlite3_column_double(stmt.handle, C.int(index))))
	case C.SQLITE_TEXT:
		pointer := C.sqlite3_column_text(stmt.handle, C.int(index))
		length := C.sqlite3_column_bytes(stmt.handle, C.int(index))
		result = value{kind: stringKind, s: C.GoStringN((*C.char)(unsafe.Pointer(pointer)), length)}
	case C.SQLITE_BLOB:
		pointer := C.sqlite3_column_blob(stmt.handle, C.int(index))
		length := C.sqlite3_column_bytes(stmt.handle, C.int(index))
		result = value{kind: stringKind, s: C.GoStringN((*C.char)(unsafe.Pointer(pointer)), length)}
	default:
		result = value{kind: nullKind}
	}
	if stmt.owner.nullMode == 1 && result.kind == stringKind && result.s == "" {
		return value{kind: nullKind}
	}
	if stmt.owner.nullMode == 2 && result.kind == nullKind {
		return value{kind: stringKind, s: ""}
	}
	return result
}
func (stmt *sqliteStatement) fetch(mode int64, column int, line int) (value, error) {
	if mode != 2 && mode != 3 && mode != 4 && mode != 5 && mode != 7 {
		return value{}, &runtimeError{line, "Error", fmt.Sprintf("gophp: PDO fetch mode %d is not implemented", mode)}
	}
	if !stmt.executed || stmt.finished {
		return boolean(false), nil
	}
	result := C.sqlite3_step(stmt.handle)
	if result == C.SQLITE_DONE {
		stmt.finished = true
		return boolean(false), nil
	}
	if result != C.SQLITE_ROW {
		err := sqliteError(line, stmt.owner.handle)
		stmt.recordError(err)
		return value{}, err
	}
	if mode == 7 {
		if column < 0 {
			return value{}, &runtimeError{line, "ValueError", "Column index must be greater than or equal to 0"}
		}
		if column >= stmt.columns {
			return value{}, &runtimeError{line, "ValueError", "Invalid column index"}
		}
		return stmt.columnValue(column), nil
	}
	row := arrayValue()
	properties := &phpArray{}
	for i := 0; i < stmt.columns; i++ {
		v := stmt.columnValue(i)
		name := stmt.columnNames[i]
		if mode != 2 && mode != 5 {
			row.a = row.a.with(arrayKey{index: int64(i)}, v)
		}
		if mode != 3 {
			row.a = row.a.with(arrayKey{isString: true, text: name}, v)
		}
		properties = properties.with(arrayKey{isString: true, text: name}, v)
	}
	if mode == 5 {
		return value{kind: objectKind, o: &phpObject{class: &classDecl{name: "stdClass", kind: "class"}, properties: properties}}, nil
	}
	return row, nil
}
func (e *evaluator) invokePDOStatement(object value, name string, args []callArgument, line int) (value, error) {
	if object.kind != objectKind || object.o == nil || object.o.pdoStatement == nil {
		return value{}, &runtimeError{line, "Error", "PDOStatement is not initialized"}
	}
	stmt := object.o.pdoStatement
	switch name {
	case "errorcode", "errorinfo":
		if len(args) != 0 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDOStatement::" + name + "() expects exactly 0 arguments"}
		}
		if name == "errorcode" {
			return pdoErrorCode(stmt.errorState), nil
		}
		return pdoErrorDetails(stmt.errorState, stmt.errorInfo), nil
	case "bindvalue":
		boundArgs, present, _, err := bindCallArguments(args, []parameter{{name: "param"}, {name: "value"}, {name: "type"}}, line)
		if err != nil {
			return value{}, err
		}
		if !present[0] || !present[1] {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDOStatement::bindValue() expects at least 2 arguments"}
		}
		if len(args) > 3 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDOStatement::bindValue() expects at most 3 arguments"}
		}
		param, _, ok := e.coerceType(dereference(boundArgs[0]), "int|string", false, e.strictTypes)
		if !ok {
			return value{}, builtinTypeError(line, "PDOStatement::bindValue", 1, "param", "string|int", boundArgs[0])
		}
		key := pdoBoundKey{index: int(param.i)}
		if param.kind == stringKind {
			key.name = param.s
			if key.name == "" {
				return value{}, &runtimeError{line, "ValueError", "PDOStatement::bindValue(): Argument #1 ($param) must not be empty"}
			}
			if !strings.HasPrefix(key.name, ":") {
				key.name = ":" + key.name
			}
			key.index = 0
		} else if key.index < 1 {
			return value{}, &runtimeError{line, "ValueError", "PDOStatement::bindValue(): Argument #1 ($param) must be greater than or equal to 1"}
		}
		paramType := int64(2)
		if present[2] {
			converted, _, ok := e.coerceType(boundArgs[2], "int", false, e.strictTypes)
			if !ok {
				return value{}, builtinTypeError(line, "PDOStatement::bindValue", 3, "type", "int", boundArgs[2])
			}
			paramType = converted.i
		}
		boundValue := dereference(boundArgs[1])
		if paramType&0xff == 2 && boundValue.kind != nullKind {
			converted, err := e.toText(boundValue, line)
			if err != nil {
				return value{}, err
			}
			boundValue = value{kind: stringKind, s: converted}
		} else if paramType&0xff == 1 && boundValue.kind == boolKind {
			boundValue = integer(boundValue.i)
		} else if paramType&0xff == 5 && boundValue.kind == intKind {
			boundValue = boolean(boundValue.i != 0)
		}
		stmt.rememberBound(key, pdoBoundValue{data: boundValue, paramType: paramType})
		return boolean(true), nil
	case "execute":
		boundArgs, present, _, err := bindCallArguments(args, []parameter{{name: "params"}}, line)
		if err != nil {
			return value{}, err
		}
		if len(args) > 1 {
			return value{}, &runtimeError{line, "ArgumentCountError", "PDOStatement::execute() expects at most 1 argument"}
		}
		var params *phpArray
		if present[0] {
			provided := dereference(boundArgs[0])
			if provided.kind != nullKind && provided.kind != arrayKind {
				return value{}, builtinTypeError(line, "PDOStatement::execute", 1, "params", "?array", provided)
			}
			if provided.kind == arrayKind {
				params = provided.a
			}
		}
		if err := e.executeSQLiteStatement(stmt, params, line); err != nil {
			return e.pdoOperationFailure(stmt.owner, "PDOStatement::execute", err, line)
		}
		return boolean(true), nil
	case "fetch", "fetchcolumn":
		mode := stmt.fetchMode
		if mode == 0 {
			mode = stmt.owner.defaultFetchMode
		}
		column := 0
		if name == "fetchcolumn" {
			if len(args) > 1 {
				return value{}, &runtimeError{line, "ArgumentCountError", "PDOStatement::fetchColumn() expects at most 1 argument"}
			}
			if len(args) == 1 {
				selected, _, ok := e.coerceType(args[0].value, "int", false, e.strictTypes)
				if !ok {
					return value{}, builtinTypeError(line, "PDOStatement::fetchColumn", 1, "column", "int", args[0].value)
				}
				column = int(selected.i)
			}
			mode = 7
		} else if len(args) > 0 && args[0].value.kind == intKind {
			mode = args[0].value.i
			if mode == 0 {
				mode = stmt.fetchMode
			}
		}
		return stmt.fetch(mode, column, line)
	case "fetchall":
		mode := stmt.fetchMode
		if mode == 0 {
			mode = stmt.owner.defaultFetchMode
		}
		if len(args) > 0 && args[0].value.kind == intKind {
			mode = args[0].value.i
			if mode == 0 {
				mode = stmt.fetchMode
			}
		}
		out := arrayValue()
		for {
			item, err := stmt.fetch(mode, 0, line)
			if err != nil {
				return value{}, err
			}
			if item.kind == boolKind && item.i == 0 {
				break
			}
			out.a = out.a.with(arrayKey{index: int64(len(out.a.entries))}, item)
		}
		return out, nil
	case "rowcount":
		return integer(stmt.affected), nil
	case "columncount":
		return integer(int64(stmt.columns)), nil
	case "closecursor":
		C.sqlite3_reset(stmt.handle)
		stmt.finished = true
		return boolean(true), nil
	case "setfetchmode":
		if len(args) > 0 {
			stmt.fetchMode = args[0].value.i
			if stmt.fetchMode == 0 {
				if e.version == PHP83 {
					// PHP 8.3 resets the statement mode to FETCH_BOTH.
					stmt.fetchMode = 4
				} else {
					stmt.fetchMode = stmt.owner.defaultFetchMode
				}
			}
		}
		return boolean(true), nil
	}
	return value{}, &runtimeError{line, "Error", "Call to undefined method PDOStatement::" + name + "()"}
}
