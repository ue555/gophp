<?php
declare(strict_types=1);
class TypedBox {}
function amount(float $x): float { return $x; }
function isMissing(?TypedBox $box): bool { return $box === null; }
echo amount(3), "|", isMissing(null) ? "null" : "object", "\n";
