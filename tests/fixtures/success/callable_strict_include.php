<?php declare(strict_types=1);
require __DIR__ . '/../support/callable_defs.inc';
$callback = importedCallable(...);
try { $callback('2'); } catch (TypeError $error) { echo 'strict|'; }
echo $callback(2), '|';
$static = ImportedCallable::run(...);
echo $static(n: 3);
