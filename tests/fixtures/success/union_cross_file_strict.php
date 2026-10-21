<?php declare(strict_types=1);
require __DIR__ . '/../support/union_strict_func.inc';
try { fromUnionFile(true); } catch (TypeError $error) { echo 'TypeError'; }
echo "\n";
