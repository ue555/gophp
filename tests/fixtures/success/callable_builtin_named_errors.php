<?php
function namedBuiltinIdentity($value) { return $value; }
try {
    call_user_func_array(unknown: 1, args: [2], callback: 'namedBuiltinIdentity');
} catch (Error $e) {
    echo $e->getMessage(), '|';
}
try {
    call_user_func_array(...['callback' => 'namedBuiltinIdentity'], callback: 'namedBuiltinIdentity', args: [2]);
} catch (Error $e) {
    echo $e->getMessage();
}
