<?php
require __DIR__ . '/../support/union_strict_func.inc';
echo gettype(fromUnionFile(5)), '|', gettype(fromUnionFile('six')), '|', gettype(fromUnionFile(true)), "\n";
