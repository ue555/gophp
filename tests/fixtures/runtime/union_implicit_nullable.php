<?php
function implicitUnion(int|string $value = null): int|string|null { return $value; }
echo gettype(implicitUnion()), "\n";
