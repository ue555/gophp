<?php
require __DIR__ . '/../support/callable_defs.inc';
$builtin = strlen(...);
$imported = importedCallable(...);
$method = ImportedCallable::run(...);
echo $builtin('abc'), '|', $imported(2), '|', $method(3), '|';
echo call_user_func_array('importedCallable', [4]), '|', call_user_func(['ImportedCallable', 'run'], 5);
