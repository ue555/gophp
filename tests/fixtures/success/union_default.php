<?php
function optionalUnion(int|string $value = 'default'): int|string { return $value; }
echo optionalUnion(), '|', optionalUnion(value: 42), "\n";
